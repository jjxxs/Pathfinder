package algorithm

import (
	"fmt"
	"strings"

	"leistungsnachweis-graphiker/problem"
)

type Algorithm interface {
	Solve(adjacency problem.Adjacency, updates chan problem.Cycle)
	Stop()
}

func FromString(algorithmName string) (Algorithm, error) {
	switch alg := strings.ToLower(algorithmName); alg {
	//case "mst":
	//	return NewMst(), nil
	case "bruteforce":
		return NewBruteForce(), nil
	//case "heldkarp":
	//	return NewHeldKarp(), nil
	default:
		return nil, fmt.Errorf("algorithm not found: %s", algorithmName)
	}
}
