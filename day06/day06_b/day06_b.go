package main

import (
	"../../util"
	"fmt"
	"strings"
)

type Node struct {
	parent    *Node
	distToSAN int
}

func (n *Node) AddChild(child *Node) {
	child.parent = n
}

func (n *Node) setDistanceToSAN(dist int) {
	n.distToSAN = dist
	if n.parent != nil {
		n.parent.setDistanceToSAN(dist + 1)
	}
}

func (n *Node) getDistanceToSAN() int {
	if n.distToSAN != -1 {
		return n.distToSAN
	} else {
		return n.parent.getDistanceToSAN() + 1
	}
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

	nodes["SAN"].setDistanceToSAN(0)
	distance := nodes["YOU"].getDistanceToSAN()

	fmt.Println(distance - 2)
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
	return &Node{nil, -1}
}
