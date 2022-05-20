package types

import (
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Node struct {
	*neo4j.Node
}

// ID allows Node to satisfy the graph.Node interface
func (n *Node) ID() int64 {
	return n.Id
}

func NewNode(id int64, labels []string, props map[string]interface{}) Node {
	return Node{
		&neo4j.Node{
			Id:     id,
			Labels: labels,
			Props:  props,
		},
	}
}

func (n *Node) String() string { return fmt.Sprintf(strings.Join(n.Labels[:], ":")) }
