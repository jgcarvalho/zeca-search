package search

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/jgcarvalho/zeca-search/rules"
)

// func ReadProbRule(fn string) ProbRule {
// 	pr := make(ProbRule)
// 	f, err := os.Open(fn)
// 	if err != nil {
// 		fmt.Println("ERROR: reading rule", err)
// 		panic(err)
// 	}
//
// 	scanner := bufio.NewScanner(f)
// 	for scanner.Scan() {
// 		var ln, c, rn string
// 		var s1, s2, s3, s4 string
// 		var p1, p2, p3, p4 float64
//
// 		r := strings.NewReplacer("[", " ", "]", " ", "->", " ", "{", " ", "}", " ", ":", " ", ",", " ")
// 		fmt.Sscanf(r.Replace(scanner.Text()), "%s %s %s %s %f %s %f %s %f %s %f", &ln, &c, &rn, &s1, &p1, &s2, &p2, &s3, &p3, &s4, &p4)
// 		pr[rules.Pattern{ln, c, rn}] = Probability{s1: p1, s2: p2, s3: p3, s4: p4}
// 	}
// 	return pr
// }

func ReadProbRule(fn string) ProbRule {
	var pr ProbRule
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("ERROR: reading rule", err)
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	r := strings.NewReplacer("[", " ", "]", " ", "->", " ", "{", " ", "}", " ", ":", " ", ",", " ")
	for scanner.Scan() {
		var ln, c, rn string
		var s1, s2, s3, s4 string
		var p1, p2, p3, p4 float64

		fmt.Sscanf(r.Replace(scanner.Text()), "%s %s %s %s %f %s %f %s %f %s %f", &ln, &c, &rn, &s1, &p1, &s2, &p2, &s3, &p3, &s4, &p4)
		// fmt.Println(rules.String2State(ln), rules.String2State(c), rules.String2State(rn), p1, p2, p3, p4)
		pr[rules.String2State(ln)][rules.String2State(c)][rules.String2State(rn)][0] = p1
		pr[rules.String2State(ln)][rules.String2State(c)][rules.String2State(rn)][1] = p2
		pr[rules.String2State(ln)][rules.String2State(c)][rules.String2State(rn)][2] = p3
		pr[rules.String2State(ln)][rules.String2State(c)][rules.String2State(rn)][3] = p4
	}
	return pr
}

func (newpk *ProbRule) Update(ind *Individual, norm int) {
	// for ln := range ind.Rule {
	// 	for c := range ind.Rule[ln] {
	// 		for rn := range ind.Rule[ln][c] {
	for ln := 0; ln < rules.NumStates; ln++ {
		for c := 0; c < rules.NumStates; c++ {
			for rn := 0; c < rules.NumStates; rn++ {
				fmt.Println(ln, c, rn, ind.Rule[ln][c][rn])
				switch ind.Rule[ln][c][rn] {
				case rules.S_c, rules.S_cn, rules.S_cp, rules.S_cG, rules.S_cP, rules.S_cneg, rules.S_cpos:
					newpk[ln][c][rn][0] += 1.0 / float64(norm)
				case rules.S_h, rules.S_hn, rules.S_hp, rules.S_hG, rules.S_hP, rules.S_hneg, rules.S_hpos:
					newpk[ln][c][rn][1] += 1.0 / float64(norm)
				case rules.S_e, rules.S_en, rules.S_ep, rules.S_eG, rules.S_eP, rules.S_eneg, rules.S_epos:
					newpk[ln][c][rn][2] += 1.0 / float64(norm)
				case rules.S_init:
					newpk[ln][c][rn][3] += 1.0 / float64(norm)
					// default:
					// panic(ind.Rule[ln][c][rn])
				}
			}
		}
	}
}

// func (newpk *ProbRule) Update(ind *Individual, norm int) {
// 	for pattern, v := range *ind.Rule {
// 		(*newpk)[pattern][v] += 1.0 / float64(norm)
// 	}
// }
//

func (pk *ProbRule) Reset() {
	var newpk ProbRule
	pk = &newpk
}

// func (pk *ProbRule) Reset() {
// 	for pattern, _ := range *pk {
// 		for k, _ := range (*pk)[pattern] {
// 			(*pk)[pattern][k] = 0.0
// 		}
// 	}
// }
//

func (pk ProbRule) Copy(newpk ProbRule) {
	pk = newpk
}

// func (pk ProbRule) Copy(newPk ProbRule) {
// 	for pattern, _ := range newPk {
// 		for k, _ := range newPk[pattern] {
// 			pk[pattern][k] = newPk[pattern][k]
// 		}
// 	}
// }

func GenRule(prob Probabilities) rules.Rule {
	var rule rules.Rule
	var rnd float64
	var i, j, k uint8
	for i = 0; i < rules.NumStates; i++ {
		for j = 0; j < rules.NumStates; j++ {
			for k = 0; k < rules.NumStates; k++ {
				rnd = rand.Float64()
				for st, val := range prob.Data[i][j][k] {
					if val > rnd {
						rule[i][j][k] = rules.Transition(rules.State(j), rules.State(st))
						break
					} else {
						rnd -= val
					}
				}
			}
		}
	}
	return rule
}

// func GenRule(prob Probabilities) rules.Rule {
// 	rule := make(rules.Rule, len(prob.Data))
// 	var tstates []string
//
// 	for k, v := range prob.Data {
// 		rnd := rand.Float64()
// 		tstates = make([]string, len(v))
// 		i := 0
// 		for s, _ := range v {
// 			tstates[i] = s
// 			i++
// 		}
// 		sort.Strings(tstates)
//
// 		for _, st := range tstates {
// 			if v[st] > rnd {
// 				rule[k] = st
// 				break
// 			} else {
// 				rnd -= v[st]
// 			}
// 		}
// 	}
// 	return rule
// }

func (prob Probabilities) Save(fn string) {
	f, err := os.Create(fn)
	if err != nil {
		fmt.Println("Error writing probabilities", err)
		panic(err)
	}
	defer f.Close()

	for ln := rules.S__; ln < rules.NumStates; ln++ {
		for c := rules.S__; c < rules.NumStates; c++ {
			for rn := rules.S__; rn < rules.NumStates; rn++ {
				if prob.Data[ln][c][rn][0]+prob.Data[ln][c][rn][1]+prob.Data[ln][c][rn][2]+prob.Data[ln][c][rn][3] > 0.0 {
					f.WriteString(fmt.Sprintf("[ %s ][ %s ][ %s ] -> { _ : %.4f, * : %.4f, | : %.4f, ? : %.4f }\n",
						rules.State2String(ln), rules.State2String(c), rules.State2String(rn),
						prob.Data[ln][c][rn][0], prob.Data[ln][c][rn][1], prob.Data[ln][c][rn][2], prob.Data[ln][c][rn][3]))
				}
			}
		}
	}
}

// func (prob Probabilities) Save(fn string) {
// 	f, err := os.Create(fn)
// 	if err != nil {
// 		fmt.Println("Error writing probabilities", err)
// 		panic(err)
// 	}
// 	defer f.Close()
//
// 	pk := prob.Data
// 	var tstates []string
// 	for k, v := range pk {
// 		// result += fmt.Sprintf("[ %s ][ %s ][ %s ] -> {", k[0], k[1], k[2])
// 		f.WriteString(fmt.Sprintf("[ %s ][ %s ][ %s ] -> {", k[0], k[1], k[2]))
// 		tstates = make([]string, len(v))
// 		i := 0
// 		for s, _ := range v {
// 			tstates[i] = s
// 			i++
// 		}
// 		sort.Strings(tstates)
// 		for _, st := range tstates {
// 			f.WriteString(fmt.Sprintf(" %s : %.4f,", st, v[st]))
// 		}
// 		f.WriteString(fmt.Sprintf("}\n"))
// 	}
// }
