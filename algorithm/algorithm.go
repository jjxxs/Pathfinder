package algorithm

import (
	"errors"
	"fmt"
)

type Algorithm interface {
}

func FromString(algorithmName string) (Algorithm, error) {
	switch algorithmName {
	case "bruteforce":
		break
	default:
		return nil, errors.New(fmt.Sprintf("algorithm not found: %s", algorithmName))
	}
}
