package pkg

import (
	"fmt"
	"sort"
)

type Node struct {
	Parent  *Node
	Service *Service
	Edges   []Node
	Depth   int
}

func (n *Node) AddEdge(node Node) {
	n.Edges = append(n.Edges, node)
}

type Resolver struct {
	resolved   map[string]bool
	unresolved map[string]bool
	services   ServiceMap
}

var (
	ErrCircularDependency = fmt.Errorf("circular dependency detected")
)

func NewDependencyGraph(services ServiceMap) (Node, error) {
	root := Node{}

	// sorting for reproducibility
	var keys []string

	for key := range services {
		// continue if the service has dependents
		for _, service := range services {
			for _, require := range service.Requires {
				if require == key {
					continue
				}
			}
		}

		keys = append(keys, key)
	}

	sort.Strings(keys)

	resolver := Resolver{
		resolved:   make(map[string]bool),
		unresolved: make(map[string]bool),
		services:   services,
	}

	for _, key := range keys {
		edge, err := resolver.graphNode(root, key)
		if err != nil {
			return Node{}, err
		}

		root.AddEdge(edge)
	}

	return root, nil
}

func (r *Resolver) graphNode(parent Node, key string) (Node, error) {
	r.unresolved[key] = true

	node := Node{
		Parent:  &parent,
		Service: r.services[key],
		Depth:   parent.Depth + 1,
	}

	for _, require := range r.services[key].Requires {
		if r.unresolved[require] {
			return Node{}, ErrCircularDependency
		}

		edge, err := r.graphNode(node, require)
		if err != nil {
			return Node{}, err
		}

		node.AddEdge(edge)
	}

	r.unresolved[key] = false

	return node, nil
}

type Walker map[string]bool

func (w Walker) Walk(node Node, f func(Node)) {
	for _, edge := range node.Edges {
		if !w[edge.Service.Name] {
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
