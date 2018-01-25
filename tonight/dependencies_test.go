package tonight

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependencyTree_leaves(t *testing.T) {
	tree := dependencyTree{
		node: &Task{ID: 1},
		children: []*dependencyTree{
			{node: &Task{ID: 11}},
			{
				node: &Task{ID: 12},
				children: []*dependencyTree{
					{node: &Task{ID: 121}},
					{node: &Task{ID: 122}},
				},
			},
		},
	}

	leaves := tree.leaves()
	ids := make([]uint, len(leaves))
	for i, leaf := range leaves {
		ids[i] = leaf.node.ID
	}

	expected := []uint{11, 121, 122}
	assert.Equal(t, expected, ids)
}

func TestDependencyTree_flat(t *testing.T) {
	tree := dependencyTree{
		node: &Task{ID: 1},
		children: []*dependencyTree{
			{node: &Task{ID: 11}},
			{
				node: &Task{ID: 12},
				children: []*dependencyTree{
					{node: &Task{ID: 121}},
					{node: &Task{ID: 122}},
				},
			},
		},
	}

	flat := tree.flat()
	ids := make([]uint, len(flat))
	for i, tree := range flat {
		ids[i] = tree.node.ID
	}

	expected := []uint{1, 11, 12, 121, 122}
	assert.Equal(t, expected, ids)
}

func TestDependencyTree_traverseBottomUp(t *testing.T) {
	tree := dependencyTree{
		node: &Task{ID: 1},
		children: []*dependencyTree{
			{node: &Task{ID: 11}},
			{
				node: &Task{ID: 12},
				children: []*dependencyTree{
					{node: &Task{ID: 121}},
					{node: &Task{ID: 122}},
				},
			},
		},
	}

	ids := make([]uint, 0)
	tree.traverseBottomUp(func(t *dependencyTree) {
		ids = append(ids, t.node.ID)
	})

	expected := []uint{122, 121, 12, 11, 1}
	assert.Equal(t, expected, ids)
}

func TestDependencyTree_buildDependencyTrees(t *testing.T) {
	tasks := []Task{
		{ID: 1},
		{ID: 11, Dependencies: []Dependency{{ID: 1}}},
		{ID: 12, Dependencies: []Dependency{{ID: 1}}},
		{ID: 121, Dependencies: []Dependency{{ID: 12}}},
		{ID: 122, Dependencies: []Dependency{{ID: 12}}},
	}

	trees := buildDependencyTrees(tasks)
	ids := make(map[uint][]uint)
	for taskID, tree := range trees {
		tree.traverseBottomUp(func(t *dependencyTree) {
			ids[taskID] = append(ids[taskID], t.node.ID)
		})
	}

	expected := map[uint][]uint{
		1:   []uint{122, 121, 12, 11, 1},
		11:  []uint{11},
		12:  []uint{122, 121, 12},
		121: []uint{121},
		122: []uint{122},
	}
	assert.Equal(t, expected, ids)
}
