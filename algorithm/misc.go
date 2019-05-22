package algorithm

import (
	"math"
	"math/bits"
	"sort"
)

// uses the sieve of atkins to generate
// a list of all primes in [2, upperBound]
// returns the primes in ascending order
func GetPrimesTo(upperBound int) []int {
	sieve := make([]bool, upperBound)

	for x := 1; x*x < upperBound; x++ {
		for y := 1; y*y < upperBound; y++ {
			n := (4 * x * x) + (y * y)
			if n <= upperBound && (n%12 == 1 || n%12 == 5) {
				sieve[n] = !sieve[n]
			}

			n = (3 * x * x) + (y * y)
			if n <= upperBound && n%12 == 7 {
				sieve[n] = !sieve[n]
			}

			n = (3 * x * x) - (y * y)
			if x > y && n <= upperBound && n%12 == 11 {
				sieve[n] = !sieve[n]
			}
		}
	}

	for i := 5; i*i < upperBound; i++ {
		if sieve[i] {
			for j := i * i; j < upperBound; j += i * i {
				sieve[j] = false
			}
		}
	}

	primes := []int{2, 3}
	for i := range sieve {
		if sieve[i] {
			primes = append(primes, i)
		}
	}

	return primes
}

// A set
type Set []int

// generates all 2^n possible subsets of the specified set
// returns them in order, ascending by length and
// and ascending by elements
func PowerSet(S Set) []Set {
	cardinality := uint(math.Pow(2, float64(len(S))))
	results := make([]Set, cardinality)

	for i := uint(0); i < cardinality; i++ {

		currentSet := Set{}
		bitMask := i
		for j := 0; j < bits.Len(i); j++ {
			if bitMask%2 == 1 {
				currentSet = append(currentSet, S[j])
			}
			bitMask = bitMask >> 1
		}

		sort.Ints(currentSet)
		results[i] = currentSet
	}

	compareByLengthOrElements := func(i, j int) bool {
		if len(results[i]) == len(results[j]) {
			for k := range results[i] {
				if results[i][k] == results[j][k] {
					continue
				}

				return results[i][k] < results[j][k]
			}
		}

		return len(results[i]) < len(results[j])
	}

	sort.Slice(results, compareByLengthOrElements)

	return results
}
