package pkg

import (
	"fmt"
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

func (s ServiceMap) NewGraph() (*Node, error) {
	root := &Node{}
	unresolved := map[string]bool{}

	for serviceName := range s {
		edge, err := s.graph(root, serviceName, unresolved)
		if err != nil {
			return nil, err
		}

		root.AddEdge(edge)
	}

	return root, nil
}

func (s ServiceMap) graph(parent *Node, name string, unresolved map[string]bool) (*Node, error) {
	node := Node{
		Parent:  parent,
		Service: s[name],
		Depth:   parent.Depth + 1,
	}

	unresolved[name] = true

	for _, require := range s[name].Requires {
		if unresolved[require] {
			return nil, fmt.Errorf("circular dependency detected: %s -> %s", name, require)
		}

		edge, err := s.graph(&node, require, unresolved)
		if err != nil {
			return nil, err
		}

		node.AddEdge(edge)
	}

	unresolved[name] = false

	return &node, nil
}

func SortNodes(node *Node) [][]*Service {
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

type Services []*Service
