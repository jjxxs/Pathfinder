package algorithm

import (
	"leistungsnachweis-graphiker/problem"
	"sort"
)

type Mst struct {
	running          bool
	shortestDistance float64
	shortestCycle    problem.Cycle
}

type edge struct {
	i    int
	j    int
	dist float64
}

func NewMst() *Mst {
	return &Mst{}
}

func (a *Mst) Stop() {
	a.running = false
}

func (a *Mst) Solve(adjacency problem.Adjacency, updates chan problem.Cycle) {
	a.running = true

	// generate all edges
	edges := make([]edge, 0)

	for i := 0; i < len(adjacency); i++ {
		for j := i + 1; j < len(adjacency); j++ {
			edges = append(edges, edge{i: i, j: j, dist: adjacency[i][j]})
		}
	}

	// sort by weight
	sort.Slice(edges, func(i, j int) bool { return edges[i].dist < edges[j].dist })

	// kruskal's algorithm
	mst := make([]edge, 0)
	for len(edges) > 0 {
		candidate := edges[0]

		// test for shortestCycle
		i, j := false, false
		for _, edge := range mst {
			if edge.i == candidate.i {
				i = true
			} else if edge.j == candidate.j {
				j = true
			}
		}

		// shortestCycle detected, discard and continue
		if i && j {
			edges = edges[1:]
			continue
		}

		// otherwise add to mst
		mst = append(mst, candidate)
		edges = edges[1:]
	}

	// generate shortestCycle
	current := 0
	a.shortestCycle = problem.Cycle{current}
	visited := make([]bool, len(adjacency))
	visited[0] = true
	v := 1
	for v < len(adjacency) {
		for _, e := range mst {
			if e.i != current && e.j != current {
				continue
			}

			if e.i == current {
				if visited[e.j] {
					current = e.j
				} else {
					visited[e.j] = true
					a.shortestCycle = append(a.shortestCycle, e.j)
					v++
				}
			} else if e.j == current {
				if visited[e.i] {
					current = e.i
				} else {
					visited[e.i] = true
					a.shortestCycle = append(a.shortestCycle, e.i)
					v++
				}
			}
		}
	}

	for i := range a.shortestCycle {
		if i == len(a.shortestCycle)-1 {
			a.shortestDistance += adjacency[a.shortestCycle[i]][a.shortestCycle[0]]
		} else {
			a.shortestDistance += adjacency[a.shortestCycle[i]][a.shortestCycle[i+1]]
		}
	}

	updates <- a.shortestCycle
	close(updates)
	a.running = false
}
