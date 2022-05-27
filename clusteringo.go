package main

import (
	"fmt"

	// "github.com/Viking2012/clusteringo/types"

	"github.com/Viking2012/clusteringo/types"
	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/topo"
)

var (
	Node1 types.Node = types.Node{Id: 1, Labels: []string{"Purchase Request"}, Properties: map[string]interface{}{"DocumentNumber": 123}}
	Node2 types.Node = types.Node{Id: 2, Labels: []string{"Purchase Order"}, Properties: map[string]interface{}{"DocumentNumber": 234}}
	Node3 types.Node = types.Node{Id: 3, Labels: []string{"Invoice"}, Properties: map[string]interface{}{"DocumentNumber": 345}}
	Node4 types.Node = types.Node{Id: 4, Labels: []string{"Fake"}, Properties: map[string]interface{}{"DocumentType": "fake document"}}

	PRLine1 types.Relationship = types.Relationship{Id: 1, Start: Node1, End: Node2, Types: []string{"OrderedWith"}, Properties: map[string]interface{}{"DocumentLine": "0001"}}
	PRLine2 types.Relationship = types.Relationship{Id: 2, Start: Node1, End: Node2, Types: []string{"OrderedWith"}, Properties: map[string]interface{}{"DocumentLine": "0002"}}
	PRLine3 types.Relationship = types.Relationship{Id: 3, Start: Node1, End: Node2, Types: []string{"OrderedWith"}, Properties: map[string]interface{}{"DocumentLine": "0003"}}
	POLine1 types.Relationship = types.Relationship{Id: 4, Start: Node2, End: Node3, Types: []string{"InvoicesWith"}, Properties: map[string]interface{}{"DocumentLine": "0001"}}
	POLine2 types.Relationship = types.Relationship{Id: 5, Start: Node2, End: Node3, Types: []string{"InvoicesWith"}, Properties: map[string]interface{}{"DocumentLine": "0002"}}
	POLine3 types.Relationship = types.Relationship{Id: 6, Start: Node2, End: Node3, Types: []string{"InvoicesWith"}, Properties: map[string]interface{}{"DocumentLine": "0003"}}
)

func main() {
	dmg := types.NewDirectedMultigraph()
	dmg.AddNode(Node1)
	dmg.AddNode(Node2)
	dmg.AddNode(Node3)
	dmg.AddRelationship(PRLine1)
	dmg.AddRelationship(PRLine2)
	dmg.AddRelationship(PRLine3)
	dmg.AddRelationship(POLine1)
	dmg.AddRelationship(POLine2)
	dmg.AddRelationship(POLine3)

	s, err := dmg.Sort()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)

	u := dmg.AsUndirected()
	cc, err := topo.ConnectedComponents(u)
	if err != nil {
		panic(err)
	}
	for i, component := range cc {
		fmt.Println(i, component)
	}

	fmt.Println()
	testMulti()
	fmt.Println()

}

func testMulti() {
	g := multi.NewDirectedGraph()
	Node1 := g.NewNode()
	g.AddNode(Node1)
	Node2 := g.NewNode()
	g.AddNode(Node2)
	Node3 := g.NewNode()
	g.AddNode(Node3)
	Node4 := g.NewNode()
	g.AddNode(Node4)

	// Manually add Lines
	g.SetLine(g.NewLine(Node1, Node2)) // PR Line 1
	g.SetLine(g.NewLine(Node1, Node2)) // PR Line 2
	g.SetLine(g.NewLine(Node1, Node2)) // PR Line 3

	g.SetLine(g.NewLine(Node2, Node3)) // PO Line 1
	g.SetLine(g.NewLine(Node2, Node3)) // PO Line 2
	g.SetLine(g.NewLine(Node2, Node3)) // PO Line 3

	s, err := topo.Sort(g)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)

	undirected := multi.NewUndirectedGraph()
	undirected.AddNode(Node1)
	undirected.AddNode(Node2)
	undirected.AddNode(Node3)
	undirected.AddNode(Node4)

	// Manually add Lines
	undirected.SetLine(undirected.NewLine(Node1, Node2)) // PR Line 1
	undirected.SetLine(undirected.NewLine(Node1, Node2)) // PR Line 2
	undirected.SetLine(undirected.NewLine(Node1, Node2)) // PR Line 3

	undirected.SetLine(undirected.NewLine(Node2, Node3)) // PO Line 1
	undirected.SetLine(undirected.NewLine(Node2, Node3)) // PO Line 2
	undirected.SetLine(undirected.NewLine(Node2, Node3)) // PO Line 3

	cc := topo.ConnectedComponents(undirected)
	for i, component := range cc {
		fmt.Println(i, component)
	}
}
