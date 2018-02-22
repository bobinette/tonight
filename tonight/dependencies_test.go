package tonight

import (
	"sort"
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
		{ID: 2},
		{ID: 21, Dependencies: []Dependency{{ID: 2}}},
		{ID: 211, Dependencies: []Dependency{{ID: 21}}},
	}

	trees := buildDependencyTrees(tasks)
	sort.Sort(byRootID(trees))
	ids := make([][]uint, len(trees))
	for i, tree := range trees {
		treeIDs := make([]uint, 0)
		tree.traverseBottomUp(func(t *dependencyTree) {
			treeIDs = append(treeIDs, t.node.ID)
		})
		ids[i] = treeIDs
	}

	expected := [][]uint{
		[]uint{122, 121, 12, 11, 1},
		[]uint{211, 21, 2},
	}
	assert.Equal(t, expected, ids)
}

type byRootID []*dependencyTree

func (l byRootID) Len() int           { return len(l) }
func (l byRootID) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byRootID) Less(i, j int) bool { return l[i].node.ID < l[j].node.ID }
