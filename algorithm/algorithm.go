package algorithm

import (
	"errors"
	"fmt"
	"strings"
)

type Algorithm interface {
	String() string
}

func FromString(algorithmName string) (Algorithm, error) {
	switch alg := strings.ToLower(algorithmName); alg {
	case "bruteforce":
		return BruteForce{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("algorithm not found: %s", algorithmName))
	}
}
