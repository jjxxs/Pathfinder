package algorithm

import (
	"log"

	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/problem"
)

type BruteForce struct{}

func (b BruteForce) String() string {
	return "Bruteforce"
}

func (b *BruteForce) Solve(adjacency [][]float64, cycles chan problem.Cycles) {
	log.Printf("solving problemset with %d entries using bruteforce", len(adjacency))
}
