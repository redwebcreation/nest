package pkg

import "sort"

type Node struct {
	Parent  *Node
	Service *Service
	Edges   []Node
	Depth   int
}

func (n *Node) AddEdge(node Node) {
	n.Edges = append(n.Edges, node)
}

func NewDependencyGraph(services ServiceMap) Node {
	root := Node{}

	// sorting for reproducibility
	var keys []string
	for key := range services {
		// We're only interested in the top level services
		if services.hasDependent(key) {
			continue
		}

		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		root.AddEdge(graphNode(root, key, services))
	}

	return root
}

func graphNode(parent Node, key string, services ServiceMap) Node {
	node := Node{
		Parent:  &parent,
		Service: services[key],
		Depth:   parent.Depth + 1,
	}

	for _, require := range services[key].Requires {
		node.AddEdge(graphNode(node, require, services))
	}

	return node
}

type Walker map[string]bool

func (w Walker) Walk(node Node, f func(Node)) {
	for _, edge := range node.Edges {
		if w[edge.Service.Name] {
			continue
		}

		w.Walk(edge, f)
	}

	// Skip the root node.
	if node.Service == nil {
		return
	}

	f(node)
	w[node.Service.Name] = true
}
