package main

import (
	"github.com/bxcodec/faker"
)

type Node struct {
	name string
	parent *Node
	children []Node
}

type fakeName struct {
	DomainName string `faker:"domain_name"`
}

func CreateTree(name string) Node {
	return Node{
		name: name,
		parent:    nil,
		children: []Node{},
	}
}

func (t *Node) Populate(depth, leafs int) {
	recursivePopulate(t, 1, depth, leafs)
}

func recursivePopulate(t *Node, currDepth, maxDepth, leafs int) {
	if currDepth == maxDepth {
		return;
	}

	var children []Node

	for i := 0; i < leafs; i++ {
		fakeName := fakeName{}
		faker.FakeData(&fakeName)

		node := Node{
			name:   fakeName.DomainName,
			parent: t,
		}

		recursivePopulate(&node, currDepth + 1, maxDepth, leafs)

		children = append(children, node)
	}

	(*t).children = children
}


