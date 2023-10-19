package loggingdrain

import (
	"encoding/json"
	"strings"
)

type LogCluster struct {
	id                int64
	logTemplateTokens []string
}

type logClusterMarshalStruct struct {
	ID                int64
	LogTemplateTokens []string
}

func (cluster *LogCluster) MarshalJSON() ([]byte, error) {
	marshalStruct := logClusterMarshalStruct{
		ID:                cluster.id,
		LogTemplateTokens: cluster.logTemplateTokens,
	}
	return json.Marshal(&marshalStruct)
}

func (cluster *LogCluster) UnmarshalJSON(data []byte) error {
	var marshalStruct logClusterMarshalStruct
	err := json.Unmarshal(data, &marshalStruct)
	if err != nil {
		return err
	}
	cluster.id = marshalStruct.ID
	cluster.logTemplateTokens = marshalStruct.LogTemplateTokens
	return nil
}

func newLogCluster(id int64, templateTokens []string) *LogCluster {
	return &LogCluster{
		id:                id,
		logTemplateTokens: templateTokens,
	}
}

func (cluster *LogCluster) getTemplate() string {
	return strings.Join(cluster.logTemplateTokens, " ")
}

type treeNodeType int

const (
	tree_node_type_root treeNodeType = iota
	tree_node_type_length
	tree_node_type_token
)

type treeNode struct {
	nodeType           treeNodeType
	length             int
	tokenNodeChildren  map[string]*treeNode
	lengthNodeChildren map[int]*treeNode
	clusters           []*LogCluster
}

type treeNodeMarshalStruct struct {
	NodeType           treeNodeType
	Length             int
	TokenNodeChildren  map[string]*treeNode
	LengthNodeChildren map[int]*treeNode
	Clusters           []*LogCluster
}

func (node *treeNode) MarshalJSON() ([]byte, error) {
	marshalStruct := treeNodeMarshalStruct{
		NodeType:           node.nodeType,
		Length:             node.length,
		TokenNodeChildren:  node.tokenNodeChildren,
		LengthNodeChildren: node.lengthNodeChildren,
		Clusters:           node.clusters,
	}
	return json.Marshal(&marshalStruct)
}

func (node *treeNode) UnmarshalJSON(data []byte) error {
	var marshalStruct treeNodeMarshalStruct
	err := json.Unmarshal(data, &marshalStruct)
	if err != nil {
		return err
	}
	node.clusters = marshalStruct.Clusters
	node.length = marshalStruct.Length
	node.lengthNodeChildren = marshalStruct.LengthNodeChildren
	node.nodeType = marshalStruct.NodeType
	node.tokenNodeChildren = marshalStruct.TokenNodeChildren
	return nil
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
		nodeType:           tree_node_type_root,
		length:             0,
		tokenNodeChildren:  make(map[string]*treeNode),
		lengthNodeChildren: make(map[int]*treeNode),
		clusters:           []*LogCluster{},
	}
}

func newLengthTreeNode(length int) *treeNode {
	return &treeNode{
		nodeType:           tree_node_type_length,
		length:             length,
		tokenNodeChildren:  make(map[string]*treeNode),
		lengthNodeChildren: make(map[int]*treeNode),
		clusters:           []*LogCluster{},
	}
}

func newTokenTreeNode() *treeNode {
	return &treeNode{
		nodeType:           tree_node_type_token,
		length:             0,
		tokenNodeChildren:  make(map[string]*treeNode),
		lengthNodeChildren: make(map[int]*treeNode),
		clusters:           []*LogCluster{},
	}
}
