package loggingdrain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogCluster(t *testing.T) {
	t.Run("test tree node", func(t *testing.T) {
		nodes := newTreeNodes()
		nodes = nodes.push(&treeNode{
			nodeType: tree_node_type_internal,
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
		assert.Equal(t, tree_node_type_internal, latestNode.nodeType)
	})
}
