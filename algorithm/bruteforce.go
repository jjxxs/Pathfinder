package algorithm

import (
	"log"
	"math"
	"time"

	"leistungsnachweis-graphiker/problem"
)

type BruteForce struct {
	ShortestDistance float32 `json:"shortestDistance"`
	ShortestCycle    []int   `json:"shortestCycle"`
	Calculations     uint64  `json:"calculations"`
	running          bool
}

func NewBruteForce() *BruteForce {
	return &BruteForce{
		ShortestDistance: math.MaxFloat32,
	}
}

func (a *BruteForce) Stop() {
	a.running = false
}

//  64.099.164
// 132.215.492
func (a *BruteForce) Solve(adjacency problem.Adjacency, cycles chan problem.Cycles) {
	// set state to running
	a.running = true
	log.Printf("solving problemset with %d entries using bruteforce", len(adjacency))

	// start worker for statistics
	go a.worker()

	// slice to permute
	points := make([]int, len(adjacency))
	for i := range points {
		points[i] = i
	}

	// calculate distance for the first permutation
	var distance float32
	for i := range points {
		if i == len(points)-1 {
			distance += adjacency[points[i]][points[0]]
		} else {
			distance += adjacency[points[i]][points[i+1]]
		}
	}

	// found new shortest cycle, set properties
	shortestCycle := make([]int, len(points))
	copy(shortestCycle, points)
	a.ShortestDistance = distance
	a.ShortestCycle = shortestCycle

	// forward result to session
	cycles <- []problem.Cycle{problem.Cycle(shortestCycle)}

	// heap's algorithm
	c := make([]int, len(adjacency))
	for i := range c {
		c[i] = 0
	}

	pointLength := len(points)
	cLength := len(c)

	i := 0
	for i < cLength && a.running {
		if c[i] < i {

			// which point to swap with
			j := 0
			if i%2 != 0 {
				j = c[i]
			}

			// indices
			jLeft := j - 1
			if jLeft < 0 {
				jLeft = pointLength - 1
			}
			jRight := j + 1
			if jRight == pointLength {
				jRight = 0
			}
			iLeft := i - 1
			if iLeft < 0 {
				iLeft = pointLength - 1
			}
			iRight := i + 1
			if iRight == pointLength {
				iRight = 0
			}

			// by only re-calculating the distances of the swapped points instead of the whole route,
			// we increase performance by a factor of two. performance increase on a intel core i7-8700k is
			// up from 64.000.000 to 132.000.000 iterations per second

			// subtract distances
			distance -= adjacency[points[j]][points[jLeft]] +
				adjacency[points[j]][points[jRight]] +
				adjacency[points[i]][points[iLeft]] +
				adjacency[points[i]][points[iRight]]

			// swap i with j
			points[j], points[i] = points[i], points[j]

			// add distances
			distance += adjacency[points[j]][points[jLeft]] +
				adjacency[points[j]][points[jRight]] +
				adjacency[points[i]][points[iLeft]] +
				adjacency[points[i]][points[iRight]]

			if distance < a.ShortestDistance {
				// found new shortest cycle, set properties and forward the result
				shortestCycle := make([]int, len(points))
				copy(shortestCycle, points)
				a.ShortestDistance = distance
				a.ShortestCycle = shortestCycle
				cycles <- []problem.Cycle{problem.Cycle(shortestCycle)}
			}

			a.Calculations++
			c[i] += 1
			i = 0
		} else {
			c[i] = 0
			i++
		}
	}

	// finished, close the channel and set state
	close(cycles)
	a.running = false
}

func (a *BruteForce) worker() {
	startTime := time.Now()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for a.running {
		<-ticker.C
		cps := float64(a.Calculations) / time.Since(startTime).Seconds()
		log.Printf("Calculations per second: %d", int64(cps))
	}
}

func (a BruteForce) String() string {
	return "Bruteforce"
}
