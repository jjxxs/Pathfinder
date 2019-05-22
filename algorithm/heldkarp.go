package algorithm

import (
	"fmt"
	"leistungsnachweis-graphiker/problem"
	"math"
)

type HeldKarp struct {
	primes   []int
	powerSet PowerSet
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

	// power set
	a.powerSet = GetPowerSet(set)

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
	for _, A := range GetPowerSet(set[1:]) {

		if len(A) == 0 {
			continue
		}

		// for every i, find the shortest distance that goes from i through A to 0
		// e.g. choose a j from A that connects to i so that distance is minimal
		for _, i := range set {

			// if i is element of A, cancel
			elementOfA := false
			for _, j := range A {
				if i == j {
					elementOfA = true
					break
				}
			}

			if elementOfA {
				continue
			}

			hash := a.getHash(A)

			// otherwise try for every j in A the route from i to j through A so that the distance is minimal
			minDistance := float32(math.MaxFloat32)
			min_j := 0
			for _, j := range A {
				A_j := make([]int, len(A)-1)
				index := 0
				for _, p := range A {
					if p == j {
						continue
					}
					A_j[index] = p
					index++
				}

				aJHash := a.getHash(A_j)
				dist := adjacency[i][j] + table[j][aJHash] // A \ {j}
				if dist < minDistance {
					minDistance = dist
					min_j = j
				}
			}

			table[i][hash] = minDistance
			backtracking[i][hash] = min_j
		}
	}

	// backtracking
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
