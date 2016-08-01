package ca

import (
	"math"

	"github.com/jgcarvalho/zeca-search/rules"
)

type Config struct {
	InitState []rules.State
	EndState  []rules.State
	// 	TransStates    []string `toml:"transition-states"`
	// 	Hydrophobicity string   `toml:"hydrophobicity"`
	// 	R              int      `toml:"r"`
	Steps int `toml:"steps"`
	// Consensus int `toml:"consensus"`
	IgnoreSteps int `toml:"ignore-steps"`
}

func (conf Config) Run(rule rules.Rule) float64 {
	var init, end, previous, current []rules.State
	init = make([]rules.State, len(conf.InitState))
	end = make([]rules.State, len(conf.EndState))
	copy(init, conf.InitState)
	copy(end, conf.EndState)
	if len(init) != len(end) {
		panic("Init and End States have diffent lenghts")
	}
	previous = make([]rules.State, len(conf.InitState))
	copy(previous, init)
	current = make([]rules.State, len(init))

	occurrence := make([]int, len(init))
	hocc := make([]uint, len(init))
	eocc := make([]uint, len(init))
	occurrence[0], occurrence[len(init)-1] = conf.Steps-conf.IgnoreSteps, conf.Steps-conf.IgnoreSteps

	// set begin and end equals to # (static states)
	current[0], current[len(init)-1] = rules.S__, rules.S__

	// fmt.Println(init)
	// fmt.Println(end)
	use := false
	for i := 0; i < conf.Steps; i++ {
		if i >= conf.IgnoreSteps {
			use = true
		}
		if i%2 == 0 {
			step(&previous, &current, &init, &end, &occurrence, &hocc, &eocc, &rule, use)
			// fmt.Println(current)
		} else {
			step(&current, &previous, &init, &end, &occurrence, &hocc, &eocc, &rule, use)
			// fmt.Println(previous)
		}

		// // change
		// for c := 1; c < len(init)-1; c++ {
		// 	current[c] = rule[rules.Pattern{previous[c-1], previous[c], previous[c+1]}]
		// 	if string(current[c][0]) == "?" {
		// 		current[c] = init[c]
		// 	}
		// 	if i >= conf.IgnoreSteps {
		// 		if current[c] == end[c] || string(end[c]) == "?" {
		// 			occurrence[c]++
		// 		}
		// 	}
		// }
		// // fmt.Println(current)
		// copy(previous, current)
		// // end change
	}
	// fmt.Println(occurrence)
	// fmt.Println("SCORE:", score(occurrence, end, conf.Steps-conf.IgnoreSteps))
	return score(occurrence, end, conf.Steps-conf.IgnoreSteps)
}

func score(oc []int, end []rules.State, norm int) float64 {
	var sc, valid float64

	for i := 0; i < len(oc); i++ {
		// esclui do calculo os estados # (inicio e fim da cadeia) e ? (indeterminado)
		if end[i] != rules.S_init && end[i] != rules.S__ {
			// if string(end[i][0]) != "?" && string(end[i][0]) != "#" {
			valid += 1.0
			if oc[i] == 0 {
				sc += math.Log(0.001)
			} else {
				sc += math.Log(float64(oc[i]) / float64(norm))
			}
		}
	}
	// return sc/float64(valid)
	return 100.0 * math.Exp(sc/float64(valid))
}

func step(previous, current, init, end *[]rules.State, occurrence *[]int, hocc, eocc *[]uint, rule *rules.Rule, use bool) {
	for c := 1; c < len(*init)-1; c++ {
		(*current)[c] = (*rule)[(*previous)[c-1]][(*previous)[c]][(*previous)[c+1]]
		if (*current)[c] == rules.S_init {
			(*current)[c] = (*init)[c]
		}
		if use {
			// ocurrence doesn't look to neighbors
			if (*current)[c] == (*end)[c] || (*end)[c] == rules.S_init {
				(*occurrence)[c]++
			}
		}
	}

	// This occurrence look to neighbors *IF IN USE, COMMENT THE OCC CODE ABOVE
	// if use {
	// 	// countOcc(current, end, occurrence)
	// 	countOcc2(current, end, occurrence, hocc, eocc)
	// }
}

func countOcc2(curr, end *[]rules.State, occurrence *[]int, hocc *[]uint, eocc *[]uint) {
	// h := make([]uint, len(*curr))
	// e := make([]uint, len(*curr))
	for c := 1; c < len(*end)-1; c++ {
		if testE3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) {
			(*eocc)[c-1], (*eocc)[c], (*eocc)[c+1] = 1, 1, 1
		} else if testH3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) {
			(*hocc)[c-1], (*hocc)[c], (*hocc)[c+1] = 1, 1, 1
		}
	}
	for c := 1; c < len(*end)-1; c++ {
		switch (*end)[c] {
		case rules.S_e, rules.S_en, rules.S_ep, rules.S_eG, rules.S_eP, rules.S_eneg, rules.S_epos:
			if (*eocc)[c] == 1 {
				(*occurrence)[c]++
			}
		case rules.S_h, rules.S_hn, rules.S_hp, rules.S_hG, rules.S_hP, rules.S_hneg, rules.S_hpos:
			if (*hocc)[c] == 1 {
				(*occurrence)[c]++
			}
		case rules.S_c, rules.S_cn, rules.S_cp, rules.S_cG, rules.S_cP, rules.S_cneg, rules.S_cpos:
			if (*eocc)[c] != 1 && (*hocc)[c] != 1 {
				(*occurrence)[c]++
			}
		}
		(*eocc)[c] = 0
		(*hocc)[c] = 0
	}
}

func countOcc(curr, end *[]rules.State, occurrence *[]int) {
	var check bool
	for c := 1; c < len(*end)-1; c++ {
		switch (*end)[c] {
		case rules.S_e, rules.S_en, rules.S_ep, rules.S_eG, rules.S_eP, rules.S_eneg, rules.S_epos:
			if (c > 0) && (c < len(*end)-3) && !check {
				if testE3((*curr)[c], (*curr)[c+1], (*curr)[c+2]) {
					(*occurrence)[c]++
					break
				}
			}
			if (c > 1) && (c < len(*end)-2) && !check {
				if testE3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) {
					(*occurrence)[c]++
					break
				}
			}
			if (c > 2) && (c < len(*end)-1) && !check {
				if testE3((*curr)[c-2], (*curr)[c-1], (*curr)[c]) {
					(*occurrence)[c]++
					break
				}
			}
		case rules.S_h, rules.S_hn, rules.S_hp, rules.S_hG, rules.S_hP, rules.S_hneg, rules.S_hpos:
			if (c > 0) && (c < len(*end)-3) && !check {
				if testH3((*curr)[c], (*curr)[c+1], (*curr)[c+2]) {
					(*occurrence)[c]++
					break
				}
			}
			if (c > 1) && (c < len(*end)-2) && !check {
				if testH3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) {
					(*occurrence)[c]++
					break
				}
			}
			if (c > 2) && (c < len(*end)-1) && !check {
				if testH3((*curr)[c-2], (*curr)[c-1], (*curr)[c]) {
					(*occurrence)[c]++
					break
				}
			}
		case rules.S_c, rules.S_cn, rules.S_cp, rules.S_cG, rules.S_cP, rules.S_cneg, rules.S_cpos:
			if (c > 0) && (c < len(*end)-3) && !check {
				if testH3((*curr)[c], (*curr)[c+1], (*curr)[c+2]) || testE3((*curr)[c], (*curr)[c+1], (*curr)[c+2]) {
					break
				}
			}
			if (c > 1) && (c < len(*end)-2) && !check {
				if testH3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) || testE3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) {
					break
				}
			}
			if (c > 2) && (c < len(*end)-1) && !check {
				if testH3((*curr)[c-2], (*curr)[c-1], (*curr)[c]) || testE3((*curr)[c-2], (*curr)[c-1], (*curr)[c]) {
					break
				}
			}
			(*occurrence)[c]++
		}
	}
}

func testE3(p0, p1, p2 rules.State) bool {
	if p0 == rules.S_e || p0 == rules.S_en || p0 == rules.S_ep || p0 == rules.S_eG || p0 == rules.S_eP || p0 == rules.S_eneg || p0 == rules.S_epos {
		if p1 == rules.S_e || p1 == rules.S_en || p1 == rules.S_ep || p1 == rules.S_eG || p1 == rules.S_eP || p1 == rules.S_eneg || p1 == rules.S_epos {
			if p2 == rules.S_e || p2 == rules.S_en || p2 == rules.S_ep || p2 == rules.S_eG || p2 == rules.S_eP || p2 == rules.S_eneg || p2 == rules.S_epos {
				return true
			}
		}
	}
	return false
}

func testH3(p0, p1, p2 rules.State) bool {
	if p0 == rules.S_h || p0 == rules.S_hn || p0 == rules.S_hp || p0 == rules.S_hG || p0 == rules.S_hP || p0 == rules.S_hneg || p0 == rules.S_hpos {
		if p1 == rules.S_h || p1 == rules.S_hn || p1 == rules.S_hp || p1 == rules.S_hG || p1 == rules.S_hP || p1 == rules.S_hneg || p1 == rules.S_hpos {
			if p2 == rules.S_h || p2 == rules.S_hn || p2 == rules.S_hp || p2 == rules.S_hG || p2 == rules.S_hP || p2 == rules.S_hneg || p2 == rules.S_hpos {
				return true
			}
		}
	}
	return false
}
