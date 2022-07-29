package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"

	"golang.org/x/exp/rand"

	"github.com/Viking2012/clusteringo/types"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/community"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

var (
	AlQaida            types.Node   = types.Node{Id: 0, Labels: []string{"Al-Qaida"}}
	Egypt              types.Node   = types.Node{Id: 1, Labels: []string{"Egypt"}}
	Hamas              types.Node   = types.Node{Id: 2, Labels: []string{"Hamas"}}
	Hezbollah          types.Node   = types.Node{Id: 3, Labels: []string{"Hezbollah"}}
	Iran               types.Node   = types.Node{Id: 4, Labels: []string{"Iran"}}
	Iraq               types.Node   = types.Node{Id: 5, Labels: []string{"Iraq"}}
	ISIS               types.Node   = types.Node{Id: 6, Labels: []string{"ISIS"}}
	Israel             types.Node   = types.Node{Id: 7, Labels: []string{"Israel"}}
	PalestianAuthority types.Node   = types.Node{Id: 8, Labels: []string{"PalestianAuthority"}}
	SaudiArabia        types.Node   = types.Node{Id: 9, Labels: []string{"SaudiArabia"}}
	Syria              types.Node   = types.Node{Id: 10, Labels: []string{"Syria"}}
	Turkey             types.Node   = types.Node{Id: 11, Labels: []string{"Turkey"}}
	USA                types.Node   = types.Node{Id: 12, Labels: []string{"USA"}}
	Countries          []types.Node = []types.Node{
		AlQaida,
		Egypt,
		Hamas,
		Hezbollah,
		Iran,
		Iraq,
		ISIS,
		Israel,
		PalestianAuthority,
		SaudiArabia,
		Syria,
		Turkey,
		USA,
	}

	// Lines: -1: Enemies, 0: Complicaged, 1: Friends
	Links [13][13]float64 = [13][13]float64{
		{-0, -1, 0, -1, 0, -1, -1, -1, -1, -1, -1, -1, -1},   // Al-Qaida
		{-1, -0, -1, -1, 0, 1, -1, 1, -1, 1, -1, -1, 1},      // Egypt
		{0, -1, -0, 0, 0, 0, -1, -1, 0, 0, -1, 1, -1},        // Hamas
		{-1, -1, 0, -0, 1, 1, -1, -1, 0, -1, 1, 0, -1},       // Hezbollah
		{0, 0, 0, 1, -0, 1, -1, -1, 0, -1, 1, -1, -1},        // Iran
		{-1, 1, 0, 1, 1, -0, -1, -1, 1, -1, 1, -1, 1},        // Iraq
		{-1, -1, -1, -1, -1, -1, -0, -1, -1, -1, -1, -1, -1}, // ISIS
		{-1, 1, 0, -1, 0, -1, -1, -0, -1, 0, -1, 0, 1},       // Israel
		{-1, 0, 0, 0, 0, 1, -1, -1, -0, 1, 0, 1, 0},          // Palestine
		{-1, 1, 0, -1, -1, -1, -1, 0, 1, -0, -1, 0, 1},       // Saudi Arabia
		{-1, -1, -1, 1, 1, 1, -1, -1, 0, -1, -0, -1, -1},     // Syria
		{-1, -1, 1, 0, -1, -1, -1, 0, 1, 0, -1, -0, 0},       // Turkey
		{-1, 1, -1, -1, -1, 1, -1, 1, 0, 1, -1, 0, -0},       // United States
	}
)

func main() {
	readJson()
	// tryLayers()
}

func tryLayers() {
	friends := simple.NewUndirectedGraph()
	complications := simple.NewUndirectedGraph()
	enemies := simple.NewUndirectedGraph()
	g := simple.NewWeightedDirectedGraph(0, 0)

	for c := 0; c < 13; c++ {
		g.AddNode(Countries[c])
		friends.AddNode(Countries[c])
		complications.AddNode(Countries[c])
		enemies.AddNode(Countries[c])
	}

	for i := 0; i < 13; i++ {
		var from types.Node = Countries[i]
		for j := (i + 1); j < 13; j++ {
			var to types.Node = Countries[j]
			var weight float64 = Links[i][j]
			nextId := int64(g.Edges().Len())
			switch weight {
			case -1:
				r := types.NewRelationshipEdge(nextId, from, to, []string{"Enemy"}, map[string]any{"weight": 1})
				enemies.SetEdge(r)
				g.SetWeightedEdge(r)
			case 0:
				r := types.NewRelationshipEdge(nextId, from, to, []string{"Complication"}, map[string]any{"weight": 5})
				complications.SetEdge(r)
				g.SetWeightedEdge(r)
			case 1:
				r := types.NewRelationshipEdge(nextId, from, to, []string{"Friend"}, map[string]any{"weight": 10})
				friends.SetEdge(r)
				g.SetWeightedEdge(r)
			}
		}
	}

	src := rand.NewSource(1)
	layers, err := community.NewUndirectedLayers(friends, enemies)
	if err != nil {
		log.Fatal(err)
	}
	weights := []float64{1, -1}

	p, err := community.Profile(
		community.ModularMultiplexScore(layers, weights, true, community.WeightMultiplex, 10, src),
		true, 1e-3, 0.1, 10,
	)
	if err != nil {
		log.Fatal(err)
	}

	for i, d := range p {
		thisG := simple.NewWeightedDirectedGraph(0, 0)
		graph.CopyWeighted(thisG, g)
		comm := d.Communities()
		for j, c := range comm {
			ByID(c)
			groupNode := types.Node{
				Id:         thisG.NewNode().ID(),
				Labels:     []string{fmt.Sprintf("Community %d, Group %d", i, j)},
				Properties: map[string]any{"weight": 100},
			}
			thisG.AddNode(groupNode)

			for _, n := range c {
				thisG.SetWeightedEdge(thisG.NewWeightedEdge(n, groupNode, 100))
			}
		}
		BySliceIDs(comm)
		fmt.Printf("Low:%.2v High:%.2v Score:%v Communities:%v Q=%.3v\n", d.Low, d.High, d.Score, comm, community.QMultiplex(layers, comm, weights, []float64{d.Low}))
		genDotFile(thisG, fmt.Sprintf("Layer %d", i))
	}
}

func BySliceIDs(c [][]graph.Node) {
	sort.Slice(c, func(i, j int) bool {
		a, b := c[i], c[j]
		l := len(a)
		if len(b) < l {
			l = len(b)
		}
		for k, v := range a[:l] {
			if v.ID() < b[k].ID() {
				return true
			}
			if v.ID() > b[k].ID() {
				return false
			}
		}
		return len(a) < len(b)
	})
}

func ByID(n []graph.Node) {
	sort.Slice(n, func(i, j int) bool { return n[i].ID() < n[j].ID() })
}

func genDotFile(g graph.Graph, filename string) (err error) {
	var b []byte
	b, err = dot.Marshal(g, "clusteringo", "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile("dotFiles/"+filename+".dot", b, 0644)
	return err
}

type jsonNode struct {
	Identity   int64          `json:"identity"`
	Labels     []string       `json:"labels"`
	Properties map[string]any `json:"properties"`
}
type jsonRelationship struct {
	Identity   int64          `json:"identity"`
	Start      int64          `json:"start"`
	End        int64          `json:"end"`
	RelType    string         `json:"type"`
	Properties map[string]any `json:"properties"`
}
type jsonFile struct {
	Nodes []jsonNode         `json:"nodes"`
	Rels  []jsonRelationship `json:"rels"`
}

func readJson() {
	raw, err := os.ReadFile(path.Join("test_data", "test.json"))
	if err != nil {
		panic(err)
	}
	var js jsonFile

	if err := json.Unmarshal(raw, &js); err != nil {
		log.Fatal(err)
	}

	var nodesMap map[int64]types.Node = make(map[int64]types.Node)
	var labels map[string]int = make(map[string]int)
	var relTypes map[string][]jsonRelationship = make(map[string][]jsonRelationship)

	for _, n := range js.Nodes {
		for _, lab := range n.Labels {
			labels[lab] += 1
		}
		nodesMap[n.Identity] = types.NewNode(n.Identity, n.Labels, n.Properties)

	}
	for _, r := range js.Rels {
		relTypes[r.RelType] = append(relTypes[r.RelType], r)
	}
	for k := range relTypes {
		log.Println(k)
	}

	var graphs map[string]*simple.UndirectedGraph = make(map[string]*simple.UndirectedGraph)
	for rType, rels := range relTypes {
		var g *simple.UndirectedGraph = simple.NewUndirectedGraph()
		for _, node := range nodesMap {
			g.AddNode(node)
		}
		for _, rel := range rels {
			from := nodesMap[rel.Start]
			to := nodesMap[rel.End]
			r := types.NewRelationshipEdge(rel.Identity, from, to, []string{rel.RelType}, rel.Properties)
			g.SetEdge(r)
		}
		graphs[rType] = g
	}

	src := rand.NewSource(1)
	layers, err := community.NewUndirectedLayers(graphs["Has_Phone"], graphs["Has_Bank"], graphs["Has_Address"], graphs["Has_Email"], graphs["Has_Fax"], graphs["Has_GUID"], graphs["Is_Named"])
	if err != nil {
		log.Fatal(err)
	}
	weights := []float64{1, 1, 1, 1, 1, 1, 1}

	p, err := community.Profile(
		community.ModularMultiplexScore(layers, weights, true, community.WeightMultiplex, 10, src),
		true, 1e-3, 0.1, 10,
	)
	if err != nil {
		log.Fatal(err)
	}

	g := simple.NewWeightedDirectedGraph(0, 0)
	for i, d := range p {
		thisG := simple.NewWeightedDirectedGraph(0, 0)
		graph.CopyWeighted(thisG, g)
		comm := d.Communities()
		for j, c := range comm {
			ByID(c)
			groupNode := types.Node{
				Id:         thisG.NewNode().ID(),
				Labels:     []string{fmt.Sprintf("Community %d, Group %d", i, j)},
				Properties: map[string]any{"weight": 100},
			}
			thisG.AddNode(groupNode)

			for _, n := range c {
				thisG.SetWeightedEdge(thisG.NewWeightedEdge(n, groupNode, 100))
			}
		}
		BySliceIDs(comm)
		fmt.Printf("Low:%.2v High:%.2v Score:%v Communities:%v Q=%.3v\n", d.Low, d.High, d.Score, comm, community.QMultiplex(layers, comm, weights, []float64{d.Low}))
		genDotFile(thisG, fmt.Sprintf("Layer %d", i))
	}
}
