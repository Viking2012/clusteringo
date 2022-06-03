package main

import (

	// "github.com/Viking2012/clusteringo/types"

	"fmt"

	"github.com/Viking2012/clusteringo/types"
	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

var (
	Node1 types.Node = types.Node{Id: 1, Labels: []string{"Purchase Request"}, Properties: map[string]interface{}{"DocumentNumber": 123}}
	Node2 types.Node = types.Node{Id: 2, Labels: []string{"Purchase Order"}, Properties: map[string]interface{}{"DocumentNumber": 234}}
	Node3 types.Node = types.Node{Id: 3, Labels: []string{"Invoice"}, Properties: map[string]interface{}{"DocumentNumber": 345}}
	Node4 types.Node = types.Node{Id: 4, Labels: []string{"Fake"}, Properties: map[string]interface{}{"DocumentType": "fake document"}}

	PREdge types.RelationshipEdge = types.RelationshipEdge{Id: 0, Start: Node1, End: Node2, Types: []string{"OrderedWith"}, Properties: map[string]interface{}{"DocumentLine": "0001"}}
	POEdge types.RelationshipEdge = types.RelationshipEdge{Id: 1, Start: Node2, End: Node3, Types: []string{"InvoicedWith"}, Properties: map[string]interface{}{"DocumentLine": "0001"}}

	PRLine1 types.RelationshipLine = types.RelationshipLine{Id: 1, Start: Node1, End: Node2, Types: []string{"OrderedWith"}, Properties: map[string]interface{}{"DocumentLine": "0001"}}
	PRLine2 types.RelationshipLine = types.RelationshipLine{Id: 2, Start: Node1, End: Node2, Types: []string{"OrderedWith"}, Properties: map[string]interface{}{"DocumentLine": "0002"}}
	PRLine3 types.RelationshipLine = types.RelationshipLine{Id: 3, Start: Node1, End: Node2, Types: []string{"OrderedWith"}, Properties: map[string]interface{}{"DocumentLine": "0003"}}
	POLine1 types.RelationshipLine = types.RelationshipLine{Id: 4, Start: Node2, End: Node3, Types: []string{"InvoicesWith"}, Properties: map[string]interface{}{"DocumentLine": "0001"}}
	POLine2 types.RelationshipLine = types.RelationshipLine{Id: 5, Start: Node2, End: Node3, Types: []string{"InvoicesWith"}, Properties: map[string]interface{}{"DocumentLine": "0002"}}
	POLine3 types.RelationshipLine = types.RelationshipLine{Id: 6, Start: Node2, End: Node3, Types: []string{"InvoicesWith"}, Properties: map[string]interface{}{"DocumentLine": "0003"}}
)

func main() {
	g := simple.NewUndirectedGraph()
	g.AddNode(Node1)
	g.AddNode(Node2)
	g.AddNode(Node3)
	g.AddNode(Node4)

	g.SetEdge(PREdge)
	g.SetEdge(POEdge)
	cc := topo.ConnectedComponents(g)
	fmt.Println("Simple: Connected Components")
	for i, component := range cc {
		fmt.Println(i, component)
	}
	fmt.Println()

	dmg := multi.NewDirectedGraph()
	dmg.AddNode(Node1)
	dmg.AddNode(Node2)
	dmg.AddNode(Node3)
	dmg.AddNode(Node4)
	dmg.SetLine(PRLine1)
	dmg.SetLine(PRLine2)
	dmg.SetLine(PRLine3)
	dmg.SetLine(POLine1)
	dmg.SetLine(POLine2)
	dmg.SetLine(POLine3)

	fmt.Println("Multi: Sort")
	s, err := topo.Sort(dmg)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	fmt.Println()

	// 	u := dmg.AsUndirected()
	// 	cc, err := topo.ConnectedComponents(u)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for i, component := range cc {
	// 		fmt.Println(i, component)
	// 	}

	// 	fmt.Println()
	// 	testMulti()
	// 	fmt.Println()

	// }

	// func testMulti() {
	// 	g := multi.NewDirectedGraph()
	// 	Node1 := g.NewNode()
	// 	g.AddNode(Node1)
	// 	Node2 := g.NewNode()
	// 	g.AddNode(Node2)
	// 	Node3 := g.NewNode()
	// 	g.AddNode(Node3)
	// 	Node4 := g.NewNode()
	// 	g.AddNode(Node4)

	// 	// Manually add Lines
	// 	g.SetLine(g.NewLine(Node1, Node2)) // PR Line 1
	// 	g.SetLine(g.NewLine(Node1, Node2)) // PR Line 2
	// 	g.SetLine(g.NewLine(Node1, Node2)) // PR Line 3

	// 	g.SetLine(g.NewLine(Node2, Node3)) // PO Line 1
	// 	g.SetLine(g.NewLine(Node2, Node3)) // PO Line 2
	// 	g.SetLine(g.NewLine(Node2, Node3)) // PO Line 3

	// 	s, err := topo.Sort(g)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(s)

	// 	undirected := multi.NewUndirectedGraph()
	// 	undirected.AddNode(Node1)
	// 	undirected.AddNode(Node2)
	// 	undirected.AddNode(Node3)
	// 	undirected.AddNode(Node4)

	// 	// Manually add Lines
	// 	undirected.SetLine(undirected.NewLine(Node1, Node2)) // PR Line 1
	// 	undirected.SetLine(undirected.NewLine(Node1, Node2)) // PR Line 2
	// 	undirected.SetLine(undirected.NewLine(Node1, Node2)) // PR Line 3

	// 	undirected.SetLine(undirected.NewLine(Node2, Node3)) // PO Line 1
	// 	undirected.SetLine(undirected.NewLine(Node2, Node3)) // PO Line 2
	// 	undirected.SetLine(undirected.NewLine(Node2, Node3)) // PO Line 3

	// 	cc := topo.ConnectedComponents(undirected)
	// 	for i, component := range cc {
	// 		fmt.Println(i, component)
	// 	}
}
