package algorithm

import (
	"leistungsnachweis-graphiker/problem"
	"math"
)

type HeldKarp struct {
	running          bool
	primes           []int
	shortestDistance float64
	shortestCycle    problem.Cycle
}

func NewHeldKarp() *HeldKarp {
	return &HeldKarp{
		shortestDistance: math.MaxFloat64,
	}
}

func (a *HeldKarp) Stop() {
	a.running = false
}

func (a *HeldKarp) Solve(adjacency problem.Adjacency, updates chan problem.Cycle) {
	a.running = true

	set := make([]int, len(adjacency))
	for i := range set {
		set[i] = i
	}

	// prime table for hashing the sets
	a.primes = GetPrimesTo(len(set) * 100)

	// create table
	table := make([]map[int]float64, len(set))

	// route from i through the empty set to 0
	for i := range set {
		table[i] = make(map[int]float64)
		table[i][a.getHash(Set{})] = adjacency[i][0]
	}

	backtracking := make([]map[int]int, len(set))
	for i := range set {
		backtracking[i] = make(map[int]int)
	}

	// for every subset of S\{0}
	for _, subset := range PowerSet(set[1:]) {

		// stopped by user
		if !a.running {
			break
		}

		if len(subset) == 0 {
			continue
		}

		// for every i, find the shortest shortestDistance that goes from i through subset to 0
		// e.g. choose a j from subset that connects to i so that shortestDistance is minimal
	NextElement:
		for _, i := range set {

			// if i is element of subset, cancel
			for _, j := range subset {
				if i == j {
					continue NextElement
				}
			}

			hash := a.getHash(subset)

			//  try for every j in subset the route from i to j through subset so that the shortestDistance is minimal
			minDistance := float64(math.MaxFloat64)
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

	if !a.running {
		return
	}

	// backtracking
	backtrackingSet := make(Set, len(set)-1)
	for i := range backtrackingSet {
		backtrackingSet[i] = i + 1
	}

	a.shortestCycle = make(problem.Cycle, len(set))
	last := 0
	for i := range set {
		h := a.getHash(backtrackingSet)
		a.shortestCycle[i] = backtracking[last][h]
		last = a.shortestCycle[i]

		tmpSet := make(Set, len(backtrackingSet)-1)

		if len(tmpSet) == 0 {
			break
		}

		index := 0
		for j := range backtrackingSet {
			if backtrackingSet[j] == a.shortestCycle[i] {
				continue
			}
			tmpSet[index] = backtrackingSet[j]
			index++
		}
		backtrackingSet = tmpSet
	}

	// get shortest shortestDistance from table
	hShortestCycle := a.getHash(set[1:])
	a.shortestDistance = table[0][hShortestCycle]

	// done, write solution to channel
	updates <- a.shortestCycle
	close(updates)
	a.running = false
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
