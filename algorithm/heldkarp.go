package algorithm

import (
	"leistungsnachweis-graphiker/problem"
)

type HeldKarp struct{}

func NewHeldKarp() *HeldKarp {
	return &HeldKarp{}
}

func (a *HeldKarp) Stop() {

}

func (a *HeldKarp) Solve(adjacency problem.Adjacency, cycles chan problem.Cycles) {
	set := make([]int, len(adjacency))
	subsets := GetPowerSet(set)

}

func (a HeldKarp) String() string {
	return "Held-Karp"
}
