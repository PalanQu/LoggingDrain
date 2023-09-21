package templateminer

import "strings"

type LogCluster struct {
	id                int64
	logTemplateTokens []string
	logs              []string
}

type LogClusterDiffType int

const (
	LOG_CLUSTER_DIFF_TYPE_NEW LogClusterDiffType = iota
	LOG_CLUSTER_DIFF_TYPE_INCREASE
	LOG_CLUSTER_DIFF_TYPE_EQUAL
	LOG_CLUSTER_DIFF_TYPE_DECREASE
)

type LogClusterDiff struct {
	DiffNum  int
	DiffRate float32
	DiffType LogClusterDiffType
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
	TREE_NODE_TYPE_ROOT treeNodeType = iota
	TREE_NODE_TYPE_LENGH
	TREE_NODE_TYPE_INTERNAL
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
		nodeType:         TREE_NODE_TYPE_ROOT,
		length:           0,
		internalChildren: make(map[string]*treeNode),
		lengthChildren:   make(map[int]*treeNode),
		clusters:         []*LogCluster{},
	}
}

func newLengthTreeNode(length int) *treeNode {
	return &treeNode{
		nodeType:         TREE_NODE_TYPE_LENGH,
		length:           length,
		internalChildren: make(map[string]*treeNode),
		lengthChildren:   make(map[int]*treeNode),
		clusters:         []*LogCluster{},
	}
}

func newInternalTreeNode() *treeNode {
	return &treeNode{
		nodeType:         TREE_NODE_TYPE_INTERNAL,
		length:           0,
		internalChildren: make(map[string]*treeNode),
		lengthChildren:   make(map[int]*treeNode),
		clusters:         []*LogCluster{},
	}
}
