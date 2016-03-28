package search

import (
	"testing"

	"github.com/jgcarvalho/zeca-search/search"
)

func TestReadProbRule(t *testing.T) {
	fn := "/home/jgcarvalho/gocode/src/github.com/jgcarvalho/zeca-create-rule/rose.rule"
	pr := search.ReadProbRule(fn)
	for k, v := range pr {
		t.Log(k, v)
	}
}

func TestGenRule(t *testing.T) {
	p := search.ReadProbRule("/home/jgcarvalho/gocode/src/github.com/jgcarvalho/zeca-create-rule/rose.rule")
	prob := &search.Probabilities{PID: 0, Generation: 0, Data: p}
	rule := search.GenRule(*prob)
	for k, v := range rule {
		t.Log(k, v)
	}
	// search.GenRule(*prob)
}
