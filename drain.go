package loggingdrain

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
)

type SearchStrategy int

const (
	SEARCH_STRATEGY_NEVER SearchStrategy = iota
	SEARCH_STRATEGY_FALLBACK
	SEARCH_STRATEGY_ALWAYS
)

type ClusterUpdateType int

const (
	CLUSTER_UPDATE_TYPE_NONE ClusterUpdateType = iota
	CLUSTER_UPDATE_TYPE_NEW_CLUSTER
	CLUSTER_UPDATE_TYPE_UPDATE_CLUSTER
)

const (
	default_max_depth    = 4
	default_sim          = 0.4
	default_max_children = 100
	default_max_clusters = 1000
	default_wildcard_str = "[*]"
)

type drain struct {
	maxDepth    int
	sim         float32
	maxChildren int
	maxClusters int

	mu             sync.Mutex
	idToCluster    *lru.Cache[int64, *LogCluster]
	clusterCounter int64
	rootNode       *treeNode
}

type drainMarshalStruct struct {
	MaxDepth    int
	Sim         float32
	MaxChildren int
	MaxClusters int

	ClusterCounter int64
	Clusters       []*LogCluster
	RootNode       *treeNode
}

func (drain *drain) MarshalJSON() ([]byte, error) {
	clusters := []*LogCluster{}
	clusters = append(clusters, drain.idToCluster.Values()...)
	marshalStruct := drainMarshalStruct{
		MaxDepth:       drain.maxDepth,
		Sim:            drain.sim,
		MaxChildren:    drain.maxChildren,
		MaxClusters:    drain.maxClusters,
		Clusters:       clusters,
		RootNode:       drain.rootNode,
		ClusterCounter: drain.clusterCounter,
	}
	return json.Marshal(&marshalStruct)
}

func (drain *drain) UnmarshalJSON(data []byte) error {
	var marshalStruct drainMarshalStruct
	err := json.Unmarshal(data, &marshalStruct)
	if err != nil {
		return err
	}
	l, _ := lru.New[int64, *LogCluster](marshalStruct.MaxClusters)
	for _, cluster := range marshalStruct.Clusters {
		l.Add(cluster.id, cluster)
	}

	drain.clusterCounter = marshalStruct.ClusterCounter
	drain.idToCluster = l
	drain.maxChildren = marshalStruct.MaxChildren
	drain.maxClusters = marshalStruct.MaxClusters
	drain.maxDepth = marshalStruct.MaxDepth
	drain.mu = sync.Mutex{}
	drain.rootNode = marshalStruct.RootNode
	drain.sim = marshalStruct.Sim
	return nil
}

func (drain *drain) status() string {
	countStr := fmt.Sprintf("cluster count %v", drain.idToCluster.Len())

	clustersStr := []string{}
	for _, clusterKey := range drain.idToCluster.Keys() {
		cluster, _ := drain.idToCluster.Get(clusterKey)
		clustersStr = append(clustersStr, fmt.Sprintf("%s\n", cluster.getTemplate()))
	}

	status := fmt.Sprintf("%s\n%s", countStr, strings.Join(clustersStr, "\n"))
	return status
}

func (drain *drain) addLogMessage(message string) (*LogCluster, ClusterUpdateType) {
	tokens := getStringTokens(message)
	cluster := drain.treeSearch(drain.rootNode, tokens, drain.sim, false)
	if cluster == nil {
		drain.clusterCounter += 1
		id := drain.clusterCounter
		cluster = newLogCluster(id, tokens)
		drain.idToCluster.Add(id, cluster)
		drain.addSeqToPrefixTree(drain.rootNode, cluster)
		return cluster, CLUSTER_UPDATE_TYPE_NEW_CLUSTER
	}
	updatedTemplate, err := drain.updateTemplate(tokens, cluster.logTemplateTokens)
	if err != nil {
		return cluster, CLUSTER_UPDATE_TYPE_NONE
	}
	if !updatedTemplate {
		return cluster, CLUSTER_UPDATE_TYPE_NONE
	}
	drain.idToCluster.Get(cluster.id)
	return cluster, CLUSTER_UPDATE_TYPE_UPDATE_CLUSTER
}

// match log message against an already existing cluster.
// Match shall be perfect (sim_th=1.0).
// New cluster will not be created as a result of this call, nor any cluster modifications.
//
// :param content: log message to match
// :param full_search_strategy: when to perform full cluster search.
//
//	(1) "never" is the fastest, will always perform a tree search [O(log(n)] but might produce
//	false negatives (wrong mismatches) on some edge cases;
//	(2) "fallback" will perform a linear search [O(n)] among all clusters with the same token count, but only in
//	case tree search found no match.
//	It should not have false negatives, however tree-search may find a non-optimal match with
//	more wildcard parameters than necessary;
//	(3) "always" is the slowest. It will select the best match among all known clusters, by always evaluating
//	all clusters with the same token count, and selecting the cluster with perfect all token match and least
//	count of wildcard matches.
//
// :return: Matched cluster or None if no match found.
func (drain *drain) match(content string, strategy SearchStrategy) *LogCluster {
	tokens := getStringTokens(content)
	requireSim := float32(1)
	fullMatch := func() *LogCluster {
		clusters := drain.getClustersForSeqLen(len(tokens))
		cluster := drain.fastMatch(clusters, tokens, requireSim, true)
		return cluster
	}
	switch strategy {
	case SEARCH_STRATEGY_ALWAYS:
		return fullMatch()
	case SEARCH_STRATEGY_FALLBACK:
		matchCluster := drain.treeSearch(drain.rootNode, tokens, requireSim, true)
		if matchCluster != nil {
			return matchCluster
		}
		return fullMatch()
	case SEARCH_STRATEGY_NEVER:
		matchCluster := drain.treeSearch(drain.rootNode, tokens, requireSim, true)
		if matchCluster != nil {
			return matchCluster
		}
		return nil
	default:
		matchCluster := drain.treeSearch(drain.rootNode, tokens, requireSim, true)
		if matchCluster != nil {
			return matchCluster
		}
		return nil
	}
}

func (drain *drain) getMaxNodeDepth() int {
	return drain.maxDepth - 2
}

func (drain *drain) GetTotalClusterSize() int {
	return drain.idToCluster.Len()
}

func (drain *drain) treeSearch(
	rootNode *treeNode,
	tokens []string,
	requireSim float32,
	includeParams bool,
) *LogCluster {
	tokenCount := len(tokens)
	lengthNode, ok := rootNode.lengthNodeChildren[tokenCount]
	if !ok {
		return nil
	}
	if tokenCount == 0 {
		if len(lengthNode.clusters) == 0 {
			return nil
		}
		return lengthNode.clusters[0]
	}
	var currentNode = lengthNode
	currentDepth := 1
	for _, token := range tokens {
		if currentDepth >= drain.getMaxNodeDepth() {
			break
		}
		if currentDepth == tokenCount {
			break
		}
		subNode, ok := currentNode.tokenNodeChildren[token]
		if !ok {
			wildcardNode := currentNode.tokenNodeChildren[default_wildcard_str]
			currentNode = wildcardNode
		} else {
			currentNode = subNode
		}
		// no wildcard node
		if currentNode == nil {
			return nil
		}
		currentDepth += 1
	}
	return drain.fastMatch(currentNode.clusters, tokens, requireSim, includeParams)
}

func (drain *drain) updateTemplate(seq1, template []string) (bool, error) {
	updated := false
	if len(seq1) != len(template) {
		return updated, errInternalRaw(
			fmt.Sprintf("seq1 length %v not equals to template length %v", len(seq1), len(template)))
	}
	for i := range template {
		if seq1[i] != template[i] && template[i] != default_wildcard_str {
			updated = true
			template[i] = default_wildcard_str
		}
	}

	return updated, nil
}

// addSeqToPrefixTree add the logCluster into tree
//
//	TODO(qujiabao): refactor the code with a smooth logic
//
//	step1: add lengthNode
//	step2: for loop each token until maxNodeDepth or latest token
//	step3: the children of current node contains the token
//		yes -> change the current node and continue loop
//		no -> the token has number
//			yes -> wildcard node in the children of current node
//				yes -> 	action1: change the current node to wildcard node and continue loop
//				no ->	action2: create the wildcard
//					  	action1: change the current node to wildcard node and continue loop
//			no -> wildcard node in the children of current node
//				yes -> the children count of current node is less than maxChildren
//					yes ->  action3: create a new node with token
//							action4: change the current node and continue loop
//					no -> action1: change the current node to wildcard node and continue loop
//				no -> the children count of current node is less than maxChildren
//					less than ->	action3: create a new node with token,
//									action4: change the current node and continue loop
//					equals -> 	action2: create the wildcard
//								action1: change the current node to wildcard node and continue loop
//					greater than -> action1: change the current node to wildcard node and continue loop
func (drain *drain) addSeqToPrefixTree(rootNode *treeNode, cluster *LogCluster) {
	tokenCount := len(cluster.logTemplateTokens)
	lengthNode, ok := rootNode.lengthNodeChildren[tokenCount]
	if !ok {
		lengthNode = newLengthTreeNode(tokenCount)
		rootNode.lengthNodeChildren[tokenCount] = lengthNode
	}
	currentNode := lengthNode
	currentDepth := 1
	if tokenCount == 0 {
		currentNode.clusters = []*LogCluster{cluster}
	}
	for _, token := range cluster.logTemplateTokens {
		if currentDepth >= drain.getMaxNodeDepth() || currentDepth >= tokenCount {
			newClusters := []*LogCluster{}
			for _, c := range currentNode.clusters {
				if drain.idToCluster.Contains(c.id) {
					newClusters = append(newClusters, c)
				}
			}
			newClusters = append(newClusters, cluster)
			currentNode.clusters = newClusters
			break
		}
		node, containsInChildren := currentNode.tokenNodeChildren[token]
		if containsInChildren {
			currentNode = node
		} else {
			wildcardNode, hasWildcardNode := currentNode.tokenNodeChildren[default_wildcard_str]
			if stringHasNumber(token) {
				if hasWildcardNode {
					currentNode = wildcardNode
				} else {
					newNode := newTokenTreeNode()
					currentNode.tokenNodeChildren[default_wildcard_str] = newNode
					currentNode = newNode
				}
			} else {
				if hasWildcardNode {
					if len(currentNode.tokenNodeChildren) < drain.maxChildren {
						newNode := newTokenTreeNode()
						currentNode.tokenNodeChildren[token] = newNode
						currentNode = newNode
					} else {
						currentNode = currentNode.tokenNodeChildren[default_wildcard_str]
					}
				} else {
					if len(currentNode.tokenNodeChildren)+1 < default_max_children {
						newNode := newTokenTreeNode()
						currentNode.tokenNodeChildren[token] = newNode
						currentNode = newNode
					} else if len(currentNode.tokenNodeChildren)+1 == default_max_children {
						newNode := newTokenTreeNode()
						currentNode.tokenNodeChildren[default_wildcard_str] = newNode
						currentNode = newNode
					} else {
						currentNode = currentNode.tokenNodeChildren[default_wildcard_str]
					}
				}
			}
		}
		currentDepth += 1
	}
}

// fastMatch find the best match for a log message (represented as tokens) versus a list of clusters
// :param clusters: List of clusters to match against
// :param tokens: the log message, separated to tokens.
// :return: Best match cluster
func (drain *drain) fastMatch(clusters []*LogCluster, tokens []string, requireSim float32, includeParams bool) *LogCluster {
	maxSim := float32(-1)
	maxParamCount := int64(-1)
	var maxMatchCluster *LogCluster

	for _, cluster := range clusters {
		if !drain.idToCluster.Contains(cluster.id) {
			continue
		}
		sim, paramCount, err := drain.getSeqDistance(cluster.logTemplateTokens, tokens, includeParams)
		if err != nil {
			continue
		}
		if sim > maxSim || (sim == maxSim && paramCount > maxParamCount) {
			maxSim = sim
			maxParamCount = paramCount
			maxMatchCluster = cluster
		}
	}
	if maxSim >= requireSim {
		return maxMatchCluster
	}

	return nil
}

func (drain *drain) getClustersForSeqLen(length int) []*LogCluster {
	stack := newTreeNodes()
	lengthNode, ok := drain.rootNode.lengthNodeChildren[length]
	if !ok {
		return []*LogCluster{}
	}
	stack = stack.push(lengthNode)
	clusters := []*LogCluster{}
	for {
		if len(stack) == 0 {
			break
		}
		var currNode *treeNode
		stack, currNode = stack.pop()
		if len(currNode.clusters) > 0 {
			clusters = append(clusters, currNode.clusters...)
		}
		for _, child := range currNode.tokenNodeChildren {
			stack = stack.push(child)
		}
	}
	return clusters
}

func (drain *drain) getSeqDistance(seq1, seq2 []string, includeParams bool) (float32, int64, error) {
	if len(seq1) != len(seq2) {
		return 0, 0, errInternalRaw(
			fmt.Sprintf("seq1 length %v not equals to seq2 length %v", len(seq1), len(seq2)))
	}
	if len(seq1) == 0 {
		return 1, 0, nil
	}
	var simTokens int64
	var paramCount int64
	for i, token1 := range seq1 {
		token2 := seq2[i]
		if token1 == default_wildcard_str {
			paramCount += 1
			continue
		}
		if token1 == token2 {
			simTokens += 1
		}
	}
	if includeParams {
		simTokens += paramCount
	}
	retVal := float32(simTokens) / float32(len(seq1))
	return retVal, paramCount, nil
}

func newDrain(options ...drainOption) *drain {
	conf := newDrainConfig(options)
	return newDrainWithConfig(conf)
}

func newDrainWithConfig(conf drainConfig) *drain {
	maxCluster := default_max_clusters
	if conf.MaxCluster > 0 {
		maxCluster = conf.MaxCluster
	}
	l, _ := lru.New[int64, *LogCluster](maxCluster)

	return &drain{
		maxDepth:       conf.Depth,
		sim:            conf.Similarity,
		maxChildren:    conf.MaxChildren,
		maxClusters:    conf.MaxCluster,
		mu:             sync.Mutex{},
		idToCluster:    l,
		clusterCounter: 0,
		rootNode:       newRootTreeNode(),
	}
}

func withDepth(depth int) drainOption {
	return drainOptionFunc(func(conf drainConfig) drainConfig {
		conf.Depth = depth
		return conf
	})
}

func withSim(sim float32) drainOption {
	return drainOptionFunc(func(conf drainConfig) drainConfig {
		conf.Similarity = sim
		return conf
	})
}

func withMaxChildren(maxChildren int) drainOption {
	return drainOptionFunc(func(conf drainConfig) drainConfig {
		conf.MaxChildren = maxChildren
		return conf
	})
}

func withMaxClusters(maxCluster int) drainOption {
	return drainOptionFunc(func(conf drainConfig) drainConfig {
		conf.MaxCluster = maxCluster
		return conf
	})
}

type drainOption interface {
	apply(drainConfig) drainConfig
}

// apply returns a config with option(s) applied.
func (o drainOptionFunc) apply(conf drainConfig) drainConfig {
	return o(conf)
}

// drainOptionFunc applies a set of options to a config.
type drainOptionFunc func(drainConfig) drainConfig

// newDrainConfig returns a config configured with options.
func newDrainConfig(options []drainOption) drainConfig {
	conf := drainConfig{
		Depth:       default_max_depth,
		Similarity:  default_sim,
		MaxChildren: default_max_children,
		MaxCluster:  default_max_clusters,
	}
	for _, o := range options {
		conf = o.apply(conf)
	}
	return conf
}
