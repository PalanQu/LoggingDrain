package loggingdrain

import "strings"

type LogCluster struct {
	id                int64
	logTemplateTokens []string
	logs              []string
}

type logClusterDiffType int

const (
	log_cluster_diff_type_new logClusterDiffType = iota
	log_cluster_diff_type_increase
	log_cluster_diff_type_equal
	log_cluster_diff_type_decrease
)

type logClusterDiff struct {
	DiffNum  int
	DiffRate float32
	DiffType logClusterDiffType
}

func newLogCluster(id int64, templateTokens []string) *LogCluster {
	return &LogCluster{
		id:                id,
		logTemplateTokens: templateTokens,
		logs:              make([]string, 0, 100),
	}
}

func (cluster *LogCluster) appendLog(log string) {
	cluster.logs = append(cluster.logs, log)
}

func (cluster *LogCluster) getTemplate() string {
	return strings.Join(cluster.logTemplateTokens, " ")
}

type treeNodeType int

const (
	tree_node_type_root treeNodeType = iota
	tree_node_type_length
	tree_node_type_internal
)

type treeNode struct {
	nodeType         treeNodeType
	length           int
	internalChildren map[string]*treeNode
	lengthChildren   map[int]*treeNode
	clusters         []*LogCluster
}

type treeNodes []*treeNode

func (nodes treeNodes) push(node *treeNode) treeNodes {
	return append(nodes, node)
}

func (nodes treeNodes) pop() (treeNodes, *treeNode) {
	latestNode := nodes[len(nodes)-1]
	return nodes[:len(nodes)-1], latestNode
}

func newTreeNodes() treeNodes {
	return []*treeNode{}
}

func newRootTreeNode() *treeNode {
	return &treeNode{
		nodeType:         tree_node_type_root,
		length:           0,
		internalChildren: make(map[string]*treeNode),
		lengthChildren:   make(map[int]*treeNode),
		clusters:         []*LogCluster{},
	}
}

func newLengthTreeNode(length int) *treeNode {
	return &treeNode{
		nodeType:         tree_node_type_length,
		length:           length,
		internalChildren: make(map[string]*treeNode),
		lengthChildren:   make(map[int]*treeNode),
		clusters:         []*LogCluster{},
	}
}

func newInternalTreeNode() *treeNode {
	return &treeNode{
		nodeType:         tree_node_type_internal,
		length:           0,
		internalChildren: make(map[string]*treeNode),
		lengthChildren:   make(map[int]*treeNode),
		clusters:         []*LogCluster{},
	}
}
