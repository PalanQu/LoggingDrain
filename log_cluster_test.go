package loggingdrain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogCluster(t *testing.T) {
	t.Run("test tree node", func(t *testing.T) {
		nodes := newTreeNodes()
		nodes = nodes.push(&treeNode{
			nodeType: tree_node_type_token,
		})
		nodes = nodes.push(&treeNode{
			nodeType: tree_node_type_length,
		})
		nodes = nodes.push(&treeNode{
			nodeType: tree_node_type_root,
		})
		assert.Equal(t, 3, len(nodes))
		nodes, latestNode := nodes.pop()
		assert.Equal(t, tree_node_type_root, latestNode.nodeType)
		assert.Equal(t, 2, len(nodes))
		nodes, latestNode = nodes.pop()
		assert.Equal(t, tree_node_type_length, latestNode.nodeType)
		nodes, latestNode = nodes.pop()
		assert.Equal(t, tree_node_type_token, latestNode.nodeType)
	})
	t.Run("test json marshal", func(t *testing.T) {
		rootNode := newRootTreeNode()
		lengthNode1 := newLengthTreeNode(3)
		lengthNode2 := newLengthTreeNode(4)
		tokenNode1 := newTokenTreeNode()
		tokenNode2 := newTokenTreeNode()
		tokenNode1.clusters = append(tokenNode1.clusters, newLogCluster(1, []string{"a", "b"}))
		tokenNode2.clusters = append(tokenNode2.clusters, newLogCluster(2, []string{"c", "d"}))
		lengthNode1.tokenNodeChildren["t1"] = tokenNode1
		lengthNode2.tokenNodeChildren["t2"] = tokenNode2
		rootNode.lengthNodeChildren[3] = lengthNode1
		rootNode.lengthNodeChildren[4] = lengthNode2
		b, err := json.Marshal(rootNode)
		if err != nil {
			t.Fatal(err)
		}

		newRootNode := newRootTreeNode()
		if err := json.Unmarshal(b, newRootNode); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, rootNode, newRootNode)
	})
}
