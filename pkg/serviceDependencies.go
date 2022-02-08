package pkg

import (
	"fmt"
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

func (n Node) String() string {
	if n.Service == nil {
		return "root"
	}

	return n.Service.Name
}

func (n *Node) Walk(f func(n *Node)) {
	if n.Service != nil {
		f(n)
	}

	for _, edge := range n.Edges {
		edge.Walk(f)
	}
}

func (s ServiceMap) NewGraph() (*Node, error) {
	graph := Graph{
		unresolved: map[string]bool{},
	}
	root := &Node{}

	for serviceName := range s {
		edge, err := graph.graph(root, serviceName, s)
		if err != nil {
			return nil, err
		}

		root.AddEdge(edge)
	}

	return root, nil
}

type Graph struct {
	unresolved map[string]bool
}

func (g Graph) graph(parent *Node, name string, services ServiceMap) (*Node, error) {
	node := Node{
		Parent:  parent,
		Service: services[name],
		Depth:   parent.Depth + 1,
	}

	g.unresolved[name] = true

	for _, require := range services[name].Requires {
		if g.unresolved[require] {
			return nil, fmt.Errorf("circular dependency detected: %s -> %s", name, require)
		}

		edge, err := g.graph(&node, require, services)
		if err != nil {
			return nil, err
		}

		node.AddEdge(edge)
	}

	g.unresolved[name] = false

	return &node, nil
}

func SortNodes(node *Node, services ServiceMap) [][]*Service {
	nodeDepth := map[string]int{}

	node.Walk(func(n *Node) {
		if nodeDepth[n.Service.Name] < n.Depth {
			nodeDepth[n.Service.Name] = n.Depth
		}
	})

	depthForNodes := map[int][]*Service{}

	for name, depth := range nodeDepth {
		depthForNodes[depth] = append(depthForNodes[depth], services[name])
	}

	sorted := make([][]*Service, len(depthForNodes))

	for key, nodes := range depthForNodes {
		sorted[len(depthForNodes)-key] = nodes
	}

	return sorted
}
