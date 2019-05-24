package algorithm

import (
	"fmt"
	"leistungsnachweis-graphiker/problem"
	"log"
	"math"
)

type HeldKarp struct {
	primes []int
}

func NewHeldKarp() *HeldKarp {
	return &HeldKarp{}
}

func (a *HeldKarp) Stop() {

}

func (a *HeldKarp) Solve(adjacency problem.Adjacency, cycles chan problem.Cycles) {
	set := make([]int, len(adjacency))
	for i := range set {
		set[i] = i
	}

	// prime table for hashing the sets
	a.primes = GetPrimesTo(len(set) * 100)

	// create table
	table := make([]map[int]float32, len(set))

	// route from i through the empty set to 0
	for i := range set {
		table[i] = make(map[int]float32)
		table[i][a.getHash(Set{})] = adjacency[i][0]
	}

	backtracking := make([]map[int]int, len(set))
	for i := range set {
		backtracking[i] = make(map[int]int)
	}

	// for every subset of S\{0}
	for _, subset := range PowerSet(set[1:]) {

		if len(subset) == 0 {
			continue
		}

		// for every i, find the shortest distance that goes from i through subset to 0
		// e.g. choose a j from subset that connects to i so that distance is minimal
	NextElement:
		for _, i := range set {

			// if i is element of subset, cancel
			for _, j := range subset {
				if i == j {
					continue NextElement
				}
			}

			hash := a.getHash(subset)

			//  try for every j in subset the route from i to j through subset so that the distance is minimal
			minDistance := float32(math.MaxFloat32)
			minJ := 0
			for _, j := range subset {

				// subset\{j}
				subsetNoJ := make([]int, len(subset)-1)
				index := 0
				for _, p := range subset {
					if p == j {
						continue
					}
					subsetNoJ[index] = p
					index++
				}
				h := a.getHash(subsetNoJ)

				dist := adjacency[i][j] + table[j][h]
				if dist < minDistance {
					minDistance = dist
					minJ = j
				}
			}

			table[i][hash] = minDistance
			backtracking[i][hash] = minJ
		}
	}

	// backtracking
	backtrackingSet := make(Set, len(set)-1)
	for i := range backtrackingSet {
		backtrackingSet[i] = i + 1
	}

	cycle := make(problem.Cycle, len(set))
	last := 0
	for i := range set {
		h := a.getHash(backtrackingSet)
		cycle[i] = backtracking[last][h]
		last = cycle[i]

		tmpSet := make(Set, len(backtrackingSet)-1)

		if len(tmpSet) == 0 {
			break
		}

		index := 0
		for j := range backtrackingSet {
			if backtrackingSet[j] == cycle[i] {
				continue
			}
			tmpSet[index] = backtrackingSet[j]
			index++
		}
		backtrackingSet = tmpSet
	}

	for _, c := range cycle {
		log.Printf("%d -> ", c)
	}

	totalDistance := float32(0.0)
	for i := range cycle {
		if i == len(cycle)-1 {
			totalDistance += adjacency[cycle[i]][0]
		} else {
			totalDistance += adjacency[cycle[i]][cycle[i+1]]
		}
	}

	log.Printf("distance: %f", totalDistance)

	s1 := backtracking[0][a.getHash(Set{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})]
	s2 := backtracking[1][a.getHash(Set{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})]
	s3 := backtracking[12][a.getHash(Set{2, 3, 4, 5, 6, 7, 8, 9, 10, 11})]
	s4 := backtracking[10][a.getHash(Set{2, 3, 4, 5, 6, 7, 8, 9, 11})]
	s5 := backtracking[7][a.getHash(Set{2, 3, 4, 5, 6, 8, 9, 11})]
	s6 := backtracking[8][a.getHash(Set{2, 3, 4, 5, 6, 9, 11})]
	s7 := backtracking[6][a.getHash(Set{2, 3, 4, 5, 9, 11})]
	s8 := backtracking[3][a.getHash(Set{2, 4, 5, 9, 11})]
	s9 := backtracking[4][a.getHash(Set{2, 5, 9, 11})]
	s10 := backtracking[5][a.getHash(Set{2, 9, 11})]
	s11 := backtracking[2][a.getHash(Set{9, 11})]
	s12 := backtracking[9][a.getHash(Set{11})]
	s13 := backtracking[11][a.getHash(Set{})]
	fmt.Printf("%d %d %d %d %d %d %d %d %d %d %d %d %d", s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12, s13)

	for x := range set {
		blaHash := a.getHash(set[1:])
		shortest := table[x][blaHash]
		fmt.Printf("done: %f", shortest)
	}
}

func (a *HeldKarp) getHash(s Set) int {
	hash := 1

	for _, k := range s {
		hash *= a.primes[k]
	}

	return hash
}

func (a HeldKarp) String() string {
	return "Held-Karp"
}
