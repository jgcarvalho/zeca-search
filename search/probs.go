package search

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/jgcarvalho/zeca-search/rules"
)

func ReadProbRule(fn string) ProbRule {
	pr := make(ProbRule)
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("ERROR: reading rule", err)
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var ln, c, rn string
		var s1, s2, s3, s4 string
		var p1, p2, p3, p4 float64

		r := strings.NewReplacer("[", " ", "]", " ", "->", " ", "{", " ", "}", " ", ":", " ", ",", " ")
		fmt.Sscanf(r.Replace(scanner.Text()), "%s %s %s %s %f %s %f %s %f %s %f", &ln, &c, &rn, &s1, &p1, &s2, &p2, &s3, &p3, &s4, &p4)
		pr[rules.Pattern{ln, c, rn}] = Probability{s1: p1, s2: p2, s3: p3, s4: p4}
	}
	return pr
}

func (pk ProbRule) Update(pop []Individual) {
	for pattern, _ := range pk {
		for k, _ := range pk[pattern] {
			pk[pattern][k] = 0.0
		}
	}

	for i := 0; i < len(pop); i++ {
		for pattern, v := range *pop[i].Rule {
			pk[pattern][v] += 1.0 / float64(len(pop))
		}
	}
}

func GenRule(prob Probabilities) rules.Rule {
	rule := make(rules.Rule, len(prob.Data))
	var tstates []string

	for k, v := range prob.Data {
		rnd := rand.Float64()
		tstates = make([]string, len(v))
		i := 0
		for s, _ := range v {
			tstates[i] = s
			i++
		}
		sort.Strings(tstates)

		for _, st := range tstates {
			if v[st] > rnd {
				rule[k] = st
				break
			} else {
				rnd -= v[st]
			}
		}
	}
	return rule
}

func (prob Probabilities) Save(fn string) {
	f, err := os.Create(fn)
	if err != nil {
		fmt.Println("Error writing probabilities", err)
		panic(err)
	}
	defer f.Close()

	pk := prob.Data
	var tstates []string
	for k, v := range pk {
		// result += fmt.Sprintf("[ %s ][ %s ][ %s ] -> {", k[0], k[1], k[2])
		f.WriteString(fmt.Sprintf("[ %s ][ %s ][ %s ] -> {", k[0], k[1], k[2]))
		tstates = make([]string, len(v))
		i := 0
		for s, _ := range v {
			tstates[i] = s
			i++
		}
		sort.Strings(tstates)
		for _, st := range tstates {
			f.WriteString(fmt.Sprintf(" %s : %.4f,", st, v[st]))
		}
		f.WriteString(fmt.Sprintf("}\n"))
	}
}
