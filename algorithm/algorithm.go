package algorithm

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"leistungsnachweis-graphiker/problem"
)

type Algorithm interface {
	String() string
	Solve(adjacency [][]float32, cycles chan problem.Cycles)
}

func FromString(algorithmName string) (Algorithm, error) {
	switch alg := strings.ToLower(algorithmName); alg {
	case "bruteforce":
		return &BruteForce{shortestDistance: math.MaxFloat32}, nil
	default:
		return nil, errors.New(fmt.Sprintf("algorithm not found: %s", algorithmName))
	}
}
