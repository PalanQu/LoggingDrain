package templateminer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogCluster(t *testing.T) {
	t.Run("test tree node", func(t *testing.T) {
		nodes := newTreeNodes()
		nodes = nodes.push(&treeNode{
			nodeType: TREE_NODE_TYPE_INTERNAL,
		})
		nodes = nodes.push(&treeNode{
			nodeType: TREE_NODE_TYPE_LENGH,
		})
		nodes = nodes.push(&treeNode{
			nodeType: TREE_NODE_TYPE_ROOT,
		})
		assert.Equal(t, 3, len(nodes))
		nodes, latestNode := nodes.pop()
		assert.Equal(t, TREE_NODE_TYPE_ROOT, latestNode.nodeType)
		assert.Equal(t, 2, len(nodes))
		nodes, latestNode = nodes.pop()
		assert.Equal(t, TREE_NODE_TYPE_LENGH, latestNode.nodeType)
		nodes, latestNode = nodes.pop()
		assert.Equal(t, TREE_NODE_TYPE_INTERNAL, latestNode.nodeType)
	})
}
