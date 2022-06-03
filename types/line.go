package types

import (
	"strings"

	"gonum.org/v1/gonum/graph"
)

var (
	EmptyRelationshipLine  RelationshipLine   = RelationshipLine{}
	EmptyRelationshipLines []RelationshipLine = []RelationshipLine{}
)

type RelationshipLine struct {
	Id         int64
	Start      Node
	End        Node
	Types      []string
	Properties map[string]any
}

func lineFactory(r RelationshipLine) graph.Line {
	return RelationshipLine{
		Id:         r.Id,
		Start:      r.Start,
		End:        r.End,
		Types:      r.Types,
		Properties: r.Properties,
	}
}

// From allows a neo4j Relationship to satisfy the interface requirements of a gonum graph.Line
func (r RelationshipLine) From() graph.Node {
	return r.Start
}

// To allows a neo4j Relationship to satisfy the interface requirements of a gonum graph.Line
func (r RelationshipLine) To() graph.Node {
	return r.End
}

// ReversedLine allows a neo4j Relationship to satisfy the interface requirements of a gonum graph.Line
func (r RelationshipLine) ReversedLine() graph.Line {
	return lineFactory(RelationshipLine{
		Id:         -r.Id,
		Start:      r.End,
		End:        r.Start,
		Types:      r.Types,
		Properties: r.Properties,
	})
}

// ID allows a neo4j Relationship to satisfy the interface requirements of a gonum graph.Line
func (r RelationshipLine) ID() int64 {
	return r.Id
}

// String prints all types of a Relationship in neo4j format
func (r *RelationshipLine) String() string { return strings.Join(r.Types[:], ":") }
