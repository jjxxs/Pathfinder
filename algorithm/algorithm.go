package algorithm

import (
	"fmt"
	"strings"

	"leistungsnachweis-graphiker/problem"
)

type Algorithm interface {
	String() string
	Solve(adjacency [][]float32, cycles chan problem.Cycles)
	Stop()
}

func FromString(algorithmName string) (Algorithm, error) {
	switch alg := strings.ToLower(algorithmName); alg {
	case "bruteforce":
		return NewBruteForce(), nil
	default:
		return nil, fmt.Errorf("algorithm not found: %s", algorithmName)
	}
}
