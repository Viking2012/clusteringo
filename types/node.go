package types

import (
	"strconv"
	"strings"
)

var (
	EmptyNode  Node   = Node{}
	EmptyNodes []Node = []Node{}
)

// Node is a representation of a neo4j driver Node
type Node struct {
	Id         int64
	Labels     []string
	Properties map[string]interface{}
}

// ID allows Node to satisfy the interface requirements of a gonum graph.Node
func (n Node) ID() int64 {
	return n.Id
}

// String prints all labels of a Node in a neo4j format
func (n Node) String() string { return strconv.Itoa(int(n.Id)) + strings.Join(n.Labels[:], ":") }
