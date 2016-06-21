package search

import "github.com/jgcarvalho/zeca-search/rules"

type Tournament []Individual

type Individual struct {
	Generation int
	Rule       *rules.Rule
	// Fitness    float64
	// Q3         float64
	Score float64
}

// type Probability map[string]float64
//
// type ProbRule map[rules.Pattern]Probability
type Probability [4]float64
type ProbRule [rules.NumStates][rules.NumStates][rules.NumStates]Probability

// type ProbRule map[[3]string]Probability

type Probabilities struct {
	// PID        uint32
	Generation int
	Data       ProbRule
}
