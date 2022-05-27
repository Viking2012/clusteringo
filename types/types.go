package types

import (
	"strconv"
	"strings"

	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/topo"
)

var (
	EmptyNode          Node           = Node{}
	EmptyRelationship  Relationship   = Relationship{}
	EmptyNodes         []Node         = []Node{}
	EmptyRelationships []Relationship = []Relationship{}
)

// Node is a representation of a neo4j driver Node
type Node struct {
	Id         int64
	Labels     []string
	Properties map[string]interface{}
}

// String prints all labels of a Node in a neo4j format
func (n Node) String() string { return strconv.Itoa(int(n.Id)) + strings.Join(n.Labels[:], ":") }

// Relationship is a representation of a neo4j driver Relationship
type Relationship struct {
	Id         int64
	Start      Node
	End        Node
	Types      []string
	Properties map[string]interface{}
}

// String prints all types of a Relationship in neo4j format
func (r *Relationship) String() string { return strings.Join(r.Types[:], ":") }

// DirectedMultiGraph is an extension of a gonum multi.DirectedGraph
// whose function arguments and returns are neo4j-type Nodes
type DirectedMultiGraph struct {
	nodes         map[int64]Node
	relationships map[int64]Relationship
	g             *multi.DirectedGraph
}

// NewDirectedMultiGraph initiates a new DirectedMultiGraph
func NewDirectedMultigraph() *DirectedMultiGraph {
	g := multi.NewDirectedGraph()
	return &DirectedMultiGraph{
		nodes:         make(map[int64]Node),
		relationships: make(map[int64]Relationship),
		g:             g,
	}
}

// Node returns a Node with the given id
func (dmg *DirectedMultiGraph) Node(id int64) Node {
	return dmg.nodes[dmg.g.Node(id).ID()]
}

// Nodes returns all nodes contained within the current graph
func (dmg *DirectedMultiGraph) Nodes() (nodes []Node) {
	for _, u := range dmg.nodes {
		nodes = append(nodes, u)
	}
	return nodes
}

// AddNode adds a new node to a DirectedMultiGraph
func (dmg *DirectedMultiGraph) AddNode(n Node) {
	new := multi.Node(n.Id)
	dmg.nodes[n.Id] = n
	dmg.g.AddNode(new)
}

// RemoveNode removes a node from the Graph, as well as all
// relationships which start or end at that node
func (dmg *DirectedMultiGraph) RemoveNode(n Node) {
	if _, ok := dmg.nodes[n.Id]; !ok {
		return
	}
	dmg.removeFroms(n)
	dmg.removeTos(n)
	dmg.g.RemoveNode(n.Id)
	delete(dmg.nodes, n.Id)
}

func (dmg *DirectedMultiGraph) removeFroms(n Node) {
	for _, f := range dmg.From(n) {
		for _, r := range dmg.RelationshipsBetween(n, f) {
			dmg.RemoveRelationship(r)
		}
	}
}

func (dmg *DirectedMultiGraph) removeTos(n Node) {
	for _, t := range dmg.To(n) {
		for _, r := range dmg.RelationshipsBetween(t, n) {
			dmg.RemoveRelationship(r)
		}
	}
}

// From returns all Nodes from which a relationship exists that Starts at n
func (dmg *DirectedMultiGraph) From(n Node) (nodes []Node) {
	i := dmg.g.From(n.Id)
	if i.Len() == 0 {
		return EmptyNodes
	}
	for i.Next() {
		u := i.Node().ID()
		nodes = append(nodes, dmg.nodes[u])
	}
	return nodes
}

// To returns all Node from which a relationship Ends at n
func (dmg *DirectedMultiGraph) To(n Node) (nodes []Node) {
	i := dmg.g.To(n.Id)
	if i.Len() == 0 {
		return EmptyNodes
	}
	for i.Next() {
		u := i.Node().ID()
		nodes = append(nodes, dmg.nodes[u])
	}
	return nodes
}

// AddRelationship
func (dmg *DirectedMultiGraph) AddRelationship(r Relationship) {
	var gLine multi.Line = multi.Line{
		F:   multi.Node(r.Start.Id),
		T:   multi.Node(r.End.Id),
		UID: r.Id,
	}
	dmg.relationships[r.Id] = r
	dmg.g.SetLine(gLine)
}

func (dmg *DirectedMultiGraph) RelationshipsBetween(u, v Node) (rels []Relationship) {
	i := dmg.g.Lines(u.Id, v.Id)
	if i.Len() == 0 {
		return EmptyRelationships
	}
	for i.Next() {
		n := i.Line().ID()
		rels = append(rels, dmg.relationships[n])
	}
	return rels
}

func (dmg *DirectedMultiGraph) Relationships() (rels []Relationship) {
	for _, r := range dmg.relationships {
		rels = append(rels, r)
	}
	return rels
}

func (dmg *DirectedMultiGraph) RemoveRelationship(r Relationship) {
	if _, ok := dmg.relationships[r.Id]; !ok {
		return
	}
	dmg.g.RemoveLine(r.Start.Id, r.End.Id, r.Id)
	delete(dmg.relationships, r.Id)
}

func (dmg *DirectedMultiGraph) Sort() (nodes []Node, err error) {
	sorted, err := topo.Sort(dmg.g)
	if err != nil {
		return EmptyNodes, err
	}

	for _, s := range sorted {
		nodes = append(nodes, dmg.nodes[s.ID()])
	}

	return nodes, nil
}

type UndirectedMultiGraph struct {
	nodes         map[int64]Node
	relationships map[int64]Relationship
	g             *multi.UndirectedGraph
}

func (dmg *DirectedMultiGraph) AsUndirected() *UndirectedMultiGraph {
	umg := NewUndirectedMultigraph()

	for _, n := range dmg.nodes {
		umg.AddNode(n)
	}

	for _, r := range dmg.relationships {
		umg.AddRelationship(r)
	}

	return umg
}

func NewUndirectedMultigraph() *UndirectedMultiGraph {
	g := multi.NewUndirectedGraph()
	return &UndirectedMultiGraph{
		nodes:         make(map[int64]Node),
		relationships: make(map[int64]Relationship),
		g:             g,
	}
}

func (umg *UndirectedMultiGraph) Node(id int64) Node {
	return umg.nodes[umg.g.Node(id).ID()]
}

func (umg *UndirectedMultiGraph) Nodes() (nodes []Node) {
	for _, u := range umg.nodes {
		nodes = append(nodes, u)
	}
	return nodes
}

func (umg *UndirectedMultiGraph) AddNode(n Node) {
	new := multi.Node(n.Id)
	umg.nodes[n.Id] = n
	umg.g.AddNode(new)
}

func (umg *UndirectedMultiGraph) RemoveNode(n Node) {
	if _, ok := umg.nodes[n.Id]; !ok {
		return
	}
	umg.removeFroms(n)
	umg.g.RemoveNode(n.Id)
	delete(umg.nodes, n.Id)
}

func (umg *UndirectedMultiGraph) removeFroms(n Node) {
	for _, f := range umg.From(n) {
		for _, r := range umg.RelationshipsBetween(n, f) {
			umg.RemoveRelationship(r)
		}
	}
}

func (umg *UndirectedMultiGraph) From(n Node) (nodes []Node) {
	i := umg.g.From(n.Id)
	if i.Len() == 0 {
		return EmptyNodes
	}
	for i.Next() {
		u := i.Node().ID()
		nodes = append(nodes, umg.nodes[u])
	}
	return nodes
}

func (umg *UndirectedMultiGraph) To(n Node) (nodes []Node) {
	i := umg.g.From(n.Id)
	if i.Len() == 0 {
		return EmptyNodes
	}
	for i.Next() {
		u := i.Node().ID()
		nodes = append(nodes, umg.nodes[u])
	}
	return nodes
}

func (umg *UndirectedMultiGraph) AddRelationship(r Relationship) {
	var gLine multi.Line = multi.Line{
		F:   multi.Node(r.Start.Id),
		T:   multi.Node(r.End.Id),
		UID: r.Id,
	}
	umg.relationships[r.Id] = r
	umg.g.SetLine(gLine)
}

func (umg *UndirectedMultiGraph) RelationshipsBetween(u, v Node) (rels []Relationship) {
	i := umg.g.Lines(u.Id, v.Id)
	if i.Len() == 0 {
		return EmptyRelationships
	}
	for i.Next() {
		n := i.Line().ID()
		rels = append(rels, umg.relationships[n])
	}
	return rels
}

func (umg *UndirectedMultiGraph) Relationships() (rels []Relationship) {
	for _, r := range umg.relationships {
		rels = append(rels, r)
	}
	return rels
}

func (umg *UndirectedMultiGraph) RemoveRelationship(r Relationship) {
	if _, ok := umg.relationships[r.Id]; !ok {
		return
	}
	umg.g.RemoveLine(r.Start.Id, r.End.Id, r.Id)
	delete(umg.relationships, r.Id)
}
