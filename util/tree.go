package util

import (
	"fmt"
	"sort"
	"strings"
)

type Node map[string]Node

func NewTree(files []string) Node {
	tree := make(Node)
	for _, file := range files {
		parts := strings.Split(file, "/")
		node := tree
		for _, part := range parts {
			if _, ok := node[part]; !ok {
				node[part] = make(Node)
			}
			node = node[part]
		}
	}
	return tree
}

func (n Node) Print(depth int) {
	if depth == 0 {
		fmt.Println("<root>")
	}
	nodeSlice := make([]string, 0, len(n))
	for key := range n {
		nodeSlice = append(nodeSlice, key)
	}

	sort.Strings(nodeSlice)

	lastObject := nodeSlice[len(nodeSlice)-1]

	for _, file := range nodeSlice {
		symbol := "├──"

		if file == lastObject {
			symbol = "└──"
		}

		fmt.Printf("%s%s %s\n", strings.Repeat("│   ", depth), symbol, file)

		if len(n[file]) != 0 {
			n[file].Print(depth + 1)
		}
	}
}
