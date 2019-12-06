package main

import (
	"../../util"
	"fmt"
	"strings"
)

type Node struct {
	depth    int
	children []*Node
}

func (n *Node) SetDepth(depth int) {
	n.depth = depth
	for _, child := range n.children {
		child.SetDepth(depth + 1)
	}
}

func (n *Node) AddChild(child *Node) {
	n.children = append(n.children, child)
}

func main() {
	nodes := make(map[string]*Node)
	util.ReadLines("day06/input.txt", func(in string) error {
		split := strings.Split(in, ")")
		parent := split[0]
		child := split[1]

		parentNode := getOrCreate(nodes, parent)
		childNode := getOrCreate(nodes, child)
		parentNode.AddChild(childNode)

		return nil
	})

	nodes["COM"].SetDepth(0)
	sum := 0
	for _, node := range nodes {
		sum += node.depth
	}
	fmt.Println(sum)
}

func getOrCreate(nodes map[string]*Node, name string) *Node {
	node, found := nodes[name]
	if found {
		return node
	} else {
		newNode := newNode()
		nodes[name] = newNode
		return newNode
	}
}

func newNode() *Node {
	return &Node{-1, make([]*Node, 0)}
}
