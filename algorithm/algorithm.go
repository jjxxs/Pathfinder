package algorithm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/problem"
)

type Algorithm interface {
	String() string
	Solve(adjacency [][]float64, cycles chan problem.Cycles)
}

func FromString(algorithmName string) (Algorithm, error) {
	switch alg := strings.ToLower(algorithmName); alg {
	case "bruteforce":
		return &BruteForce{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("algorithm not found: %s", algorithmName))
	}
}
