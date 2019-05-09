package algorithm

import (
	"log"
	"time"

	"leistungsnachweis-graphiker/problem"
)

type BruteForce struct {
	shortestDistance float32
	shortestCycle    []int
	calculations     int64
}

func (b BruteForce) String() string {
	return "Bruteforce"
}

func (b *BruteForce) Solve(adjacency [][]float32, cycles chan problem.Cycles) {
	log.Printf("solving problemset with %d entries using bruteforce", len(adjacency))

	// slice to permutate
	points := make([]int, len(adjacency))
	for i := range points {
		points[i] = i
	}

	// heap's algorithm
	c := make([]int, len(adjacency))
	for i := range c {
		c[i] = 0
	}

	go b.worker()

	i := 0
	for i < len(c) {
		if c[i] < i {
			if i%2 == 0 {
				points[0], points[i] = points[i], points[0]
			} else {
				points[c[i]], points[i] = points[i], points[c[i]]
			}

			// new permutation, calculate distance
			distance := float32(0)
			for i := range points {
				if i == len(points)-1 {
					distance += adjacency[points[i]][points[0]]
				} else {
					distance += adjacency[points[i]][points[i+1]]
				}
			}

			if distance < b.shortestDistance {
				b.shortestDistance = distance
				//shortestCycle := make([]int, len(points))
				//copy(points, shortestCycle)
				//b.shortestCycle = shortestCycle
				log.Printf("bruteforce found new shortest cycle: %f", distance)
			}

			b.calculations++

			c[i] += 1
			i = 0

		} else {
			c[i] = 0
			i++
		}
	}
}

func (b *BruteForce) worker() {
	for {
		time.Sleep(1 * time.Second)
		log.Printf("calculations per second: %d", b.calculations)
		b.calculations = 0
	}
}
