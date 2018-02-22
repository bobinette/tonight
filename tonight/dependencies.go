package tonight

type dependencyTree struct {
	node     *Task
	children []*dependencyTree
}

func buildDependencyTrees(tasks []Task) []*dependencyTree {
	// Create all the nodes
	roots := make(map[uint]struct{})
	nodes := make(map[uint]*dependencyTree)
	for i, task := range tasks {
		nodes[task.ID] = &dependencyTree{
			node: &tasks[i],
		}
		roots[task.ID] = struct{}{}
	}

	for _, task := range tasks {
		for _, dep := range task.Dependencies {
			if _, ok := nodes[dep.ID]; ok {
				nodes[dep.ID].children = append(nodes[dep.ID].children, nodes[task.ID])
				delete(roots, task.ID)
			}
		}
	}

	trees := make([]*dependencyTree, 0)
	for rootID := range roots {
		trees = append(trees, nodes[rootID])
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
