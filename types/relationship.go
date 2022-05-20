package types

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

type Relationship struct {
	*neo4j.Relationship
}

// From allows Relationship to satisfy the graph.Line interface
func (r *Relationship) From() Node {
	return Node{&neo4j.Node{}}
}

// To allows Relationship to satisfy the graph.Line interface
func (r *Relationship) To() Node {
	return Node{&neo4j.Node{}}
}

// ReversedLine allows Relationship to satisfy the graph.Line interface
func (r *Relationship) ReveresedLine() Relationship {
	newR := Relationship{
		&neo4j.Relationship{
			Id:      r.Id,
			StartId: r.EndId,   // this is the reversal
			EndId:   r.StartId, // this is the reversal
			Type:    r.Type,
			Props:   r.Props,
		},
	}
	return newR
}

// ID allows Relationship to satisfy the graph.Line interface
func (r *Relationship) ID() int64 {
	return r.Id
}
