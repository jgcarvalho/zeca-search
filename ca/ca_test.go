package ca

import (
	"testing"

	"github.com/jgcarvalho/zeca-search/ca"
	"github.com/jgcarvalho/zeca-search/search"
)

func TestRun(t *testing.T) {
	p := search.ReadProbRule("/home/jgcarvalho/gocode/src/github.com/jgcarvalho/zeca-create-rule/simple.rule")
	prob := &search.Probabilities{PID: 0, Generation: 0, Data: p}
	rule := search.GenRule(*prob)
	cellAuto := ca.Config{
		InitState:   []string{"#", "M", "A", "D", "F", "G", "H", "I", "K", "#", "A", "A", "#"},
		EndState:    []string{"#", "_", "_", "*", "?", "*", "*", "_", "|", "#", "|", "_", "#"},
		Steps:       10,
		IgnoreSteps: 2}
	cellAuto.Run(rule)
	// for k, v := range rule {
	// 	t.Log(k, v)
	// }
	// search.GenRule(*prob)
}
