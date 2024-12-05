package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	g := NewGraph()

	modA := ModulePublic{Path: "A", Version: "1.0"}
	modB := ModulePublic{Path: "B", Version: "1.1"}
	modC := ModulePublic{Path: "C", Version: "1.2"}
	modD := ModulePublic{Path: "D", Version: "1.0"}
	modE := ModulePublic{Path: "E", Version: "1.1"}
	modF := ModulePublic{Path: "F", Version: "1.2"}
	modG := ModulePublic{Path: "G", Version: "1.0"}
	modH := ModulePublic{Path: "H", Version: "1.0"}

	g.AddEdge(modA, modB)
	g.AddEdge(modB, modC)
	g.AddEdge(modC, modA) // Cycle: A -> B -> C -> A
	g.AddEdge(modD, modE)
	g.AddEdge(modE, modF)
	g.AddEdge(modF, modD) // Cycle: D -> E -> F -> D
	g.AddEdge(modC, modD) // Link between two SCCs
	g.AddEdge(modG, modA)
	g.AddEdge(modD, modH)

	sccs := g.FindSCCs()

	expected := [][]ModulePublic{
		{
			ModulePublic{Path: "H", Version: "1.0"},
		},
		{
			ModulePublic{Path: "D", Version: "1.0"},
			ModulePublic{Path: "E", Version: "1.1"},
			ModulePublic{Path: "F", Version: "1.2"},
		},
		{
			ModulePublic{Path: "A", Version: "1.0"},
			ModulePublic{Path: "B", Version: "1.1"},
			ModulePublic{Path: "C", Version: "1.2"},
		},
		{
			ModulePublic{Path: "G", Version: "1.0"},
		},
	}

	assert.Equal(t, expected, sccs)
}
