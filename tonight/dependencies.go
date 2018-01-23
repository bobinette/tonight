package tonight

type dependencyTree struct {
	node     *Task
	children []*dependencyTree
}

func buildDependencyTrees(tasks []Task) map[uint]*dependencyTree {
	trees := make(map[uint]*dependencyTree)

	// Create all the nodes
	for i, task := range tasks {
		trees[task.ID] = &dependencyTree{
			node: &tasks[i],
		}
	}

	for _, task := range tasks {
		for _, dep := range task.Dependencies {
			if _, ok := trees[dep.ID]; ok {
				trees[dep.ID].children = append(trees[dep.ID].children, trees[task.ID])
			}
		}
	}

	return trees
}

func (t *dependencyTree) leaves() []*dependencyTree {
	if len(t.children) == 0 {
		return []*dependencyTree{t}
	}

	leaves := make([]*dependencyTree, 0)
	for _, child := range t.children {
		leaves = append(leaves, child.leaves()...)
	}
	return leaves
}

func (t *dependencyTree) flat() []*dependencyTree {
	res := make([]*dependencyTree, 0)
	current := []*dependencyTree{t}

	for len(current) > 0 {
		next := make([]*dependencyTree, 0)
		for _, tree := range current {
			res = append(res, tree)
			next = append(next, tree.children...)
		}
		current = next
	}

	return res
}

func (t *dependencyTree) traverseBottomUp(f func(*dependencyTree)) {
	flat := t.flat()
	for i := len(flat) - 1; i >= 0; i-- {
		f(flat[i])
	}
}
