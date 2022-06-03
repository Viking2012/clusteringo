package types

import (
	"strings"

	"gonum.org/v1/gonum/graph"
)

var (
	EmptyRelationshipEdge  RelationshipEdge   = RelationshipEdge{}
	EmptyRelationshipEdges []RelationshipEdge = []RelationshipEdge{}
)

// RelationshipEdge is a representation of a neo4j driver RelationshipEdge
type RelationshipEdge struct {
	Id         int64
	Start      Node
	End        Node
	Types      []string
	Properties map[string]any
}

// edgeFactory is required for go to register a RelationshipEdge as a gonum graph.Edge
// Otherwise, gonum functions which receive or return a graph.Edge interfaces wouldn't recognize
// a RelationshipEdge as satisfying the interface.
func edgeFactory(r RelationshipEdge) graph.Edge {
	return RelationshipEdge{
		Id:         r.Id,
		Start:      r.End,
		End:        r.Start,
		Types:      r.Types,
		Properties: r.Properties,
	}
}

// From allows a neo4j Relationship to satisfy the interface requirements of a gonum graph.Edge
func (r RelationshipEdge) From() graph.Node {
	return r.Start
}

// To allows a neo4j Relationship to satisfy the interface requirements of a gonum graph.Edge
func (r RelationshipEdge) To() graph.Node {
	return r.End
}

// ReveresedEdge allows a neo4j Relationship to satisfy the interface requirements of a gonum graph.Edge
func (r RelationshipEdge) ReversedEdge() graph.Edge {
	return edgeFactory(r)
}

// String prints all types of a Relationship in neo4j format
func (r *RelationshipEdge) String() string { return strings.Join(r.Types[:], ":") }
