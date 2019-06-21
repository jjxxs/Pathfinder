package algorithm

import (
	"leistungsnachweis-graphiker/problem"
	"math"
	"testing"
)

func TestHeldKarp(t *testing.T) {
	points := []problem.Point{
		{X: 0, Y: 0},
		{X: 50, Y: 0},
		{X: 50, Y: 50},
		{X: 0, Y: 50},
		{X: 25, Y: 75},
	}

	p := problem.NewProblem(points)
	b := NewHeldKarp()
	u := make(chan problem.Cycle, 10)

	go b.Solve(p.Adjacency, u)

	for {
		cycle, hasMore := <-u
		if !hasMore {
			break
		}
		p.UpdateRoute(cycle)
	}

	if math.Round(p.ShortestDistance*100)/100 != 220.71 {
		t.Fatalf("wrong distance: %f", p.ShortestDistance)
	}
}
