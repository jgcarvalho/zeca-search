package search

import (
	"testing"

	"github.com/jgcarvalho/zeca-search/rules"
)

// "github.com/jgcarvalho/zeca-search/search"

func TestReadProbRule(t *testing.T) {
	fn := "/home/jgcarvalho/gocode/src/github.com/jgcarvalho/zeca-create-rule/rose.rule"
	pr := ReadProbRule(fn)
	t.Log(pr)
}

func TestGenRule(t *testing.T) {
	p := ReadProbRule("/home/jgcarvalho/gocode/src/github.com/jgcarvalho/zeca-create-rule/rose.rule")
	prob := new(Probabilities)
	prob.Generation = 0
	prob.Data = p
	rule := GenRule(*prob)
	for i := rules.S_An; i < rules.S_ep; i++ {
		for j := rules.S_An; j < rules.S_ep; j++ {
			for k := rules.S_An; k < rules.S_ep; k++ {
				t.Log(i, j, k, rule[i][j][k])
			}
		}
	}
}
