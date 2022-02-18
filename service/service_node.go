package service

import (
	"sort"
)

type Node struct {
	Parent  *Node
	Service *Service
	Edges   []*Node
	Depth   int
}

func (n *Node) AddEdge(e *Node) {
	n.Edges = append(n.Edges, e)
}

func (n *Node) Walk(f func(n *Node)) {
	// skip the root node
	if n.Service != nil {
		f(n)
	}

	for _, edge := range n.Edges {
		edge.Walk(f)
	}
}

func sortNodes(node *Node) [][]*Service {
	nodeDepth := map[*Service]int{}
	depthNode := map[int][]*Service{}

	node.Walk(func(n *Node) {
		if nodeDepth[n.Service] < n.Depth {
			nodeDepth[n.Service] = n.Depth
		}
	})

	for service, depth := range nodeDepth {
		depthNode[depth] = append(depthNode[depth], service)
	}

	reversed := make([][]*Service, len(depthNode))

	for key, nodes := range depthNode {
		// sort nodes for reproducibility
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Name < nodes[j].Name
		})

		// reverse the order of the nodes so that the node with the highest depth comes first
		reversed[len(depthNode)-key] = nodes
	}

	return reversed
}
