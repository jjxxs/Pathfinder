package algorithm

import (
	"fmt"
	"strings"

	"leistungsnachweis-graphiker/problem"
)

type Algorithm interface {
	String() string
	Solve(adjacency problem.Adjacency, cycles chan problem.Cycles)
	Stop()
	GetSolution() (error, float64, problem.Cycle)
}

func FromString(algorithmName string) (Algorithm, error) {
	switch alg := strings.ToLower(algorithmName); alg {
	case "bruteforce":
		return NewBruteForce(), nil
	case "heldkarp":
		return NewHeldKarp(), nil
	default:
		return nil, fmt.Errorf("algorithm not found: %s", algorithmName)
	}
}
