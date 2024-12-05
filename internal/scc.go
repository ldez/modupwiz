package internal

import (
	"slices"
	"strings"
)

type Graph struct {
	adjacencyList map[string][]ModulePublic
	nodes         map[string]ModulePublic
}

func NewGraph() *Graph {
	return &Graph{
		adjacencyList: make(map[string][]ModulePublic),
		nodes:         make(map[string]ModulePublic),
	}
}

func (g *Graph) AddEdge(from, to ModulePublic) {
	g.adjacencyList[from.Path] = append(g.adjacencyList[from.Path], to)
	g.nodes[from.Path] = from
	g.nodes[to.Path] = to
}

// FindSCCs Tarjan's Algorithm to find SCCs.
func (g *Graph) FindSCCs() [][]ModulePublic {
	index := 0

	var stack []string

	onStack := make(map[string]bool)
	indices := make(map[string]int)
	lowLink := make(map[string]int)

	var sccs [][]ModulePublic

	var dfs func(v string)
	dfs = func(v string) {
		indices[v] = index
		lowLink[v] = index

		index++

		stack = append(stack, v)

		onStack[v] = true

		for _, neighbor := range g.adjacencyList[v] {
			if _, ok := indices[neighbor.Path]; !ok {
				dfs(neighbor.Path)
				lowLink[v] = min(lowLink[v], lowLink[neighbor.Path])
			} else if onStack[neighbor.Path] {
				lowLink[v] = min(lowLink[v], indices[neighbor.Path])
			}
		}

		if lowLink[v] == indices[v] {
			var scc []ModulePublic

			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				onStack[w] = false

				scc = append(scc, g.nodes[w])

				if w == v {
					break
				}
			}

			slices.SortFunc(scc, func(a, b ModulePublic) int {
				return strings.Compare(a.Path, b.Path)
			})

			sccs = append(sccs, scc)
		}
	}

	for node := range g.adjacencyList {
		if _, ok := indices[node]; !ok {
			dfs(node)
		}
	}

	return sccs
}
