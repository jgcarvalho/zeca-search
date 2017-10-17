package ca

import (
	"fmt"
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
	FitFunc     string
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

	// Confusion Matrix
	var cm [3][4]int

	occurrence := make([]int, len(init))
	hocc := make([]uint, len(init))
	eocc := make([]uint, len(init))
	// occurrence[0], occurrence[len(init)-1] = conf.Steps-conf.IgnoreSteps, conf.Steps-conf.IgnoreSteps

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
			step(&previous, &current, &init, &end, &cm, &occurrence, &hocc, &eocc, &rule, use)
			// fmt.Println(current)
		} else {
			step(&current, &previous, &init, &end, &cm, &occurrence, &hocc, &eocc, &rule, use)
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
	fmt.Println("CM", cm)
	// fmt.Println("SCORE:", score(occurrence, end, conf.Steps-conf.IgnoreSteps))
	a := mcc2(cm)
	b := CBA(cm)
	// fmt.Println("Acc", score(occurrence, end, conf.Steps-conf.IgnoreSteps))
	fmt.Println("MCC", mcc(cm))
	fmt.Println("MCC2", a)
	fmt.Println("CBA", b)
	switch conf.FitFunc {
	case "mcc":
		return a
	case "cba":
		return b
	case "cba+mcc":
		return a + b
	case "cba*mcc":
		return a * b
	case "cba2*mcc":
		return a * b * b
	default:
		return a + b
	}
	// return score(occurrence, end, conf.Steps-conf.IgnoreSteps)
}

func score(oc []int, end []rules.State, norm int) float64 {
	var sc, valid float64

	for i := 0; i < len(oc); i++ {
		// exclui do calculo os estados # (inicio e fim da cadeia) e ? (indeterminado)
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

func mcc(cm [3][4]int) float64 {
	tp := float64(cm[0][0] + cm[1][1] + cm[2][2])
	tn := float64((cm[1][1] + cm[1][2] + cm[2][1] + cm[2][2]) + (cm[0][0] + cm[0][2] + cm[2][0] + cm[2][2]) + (cm[0][0] + cm[0][1] + cm[1][0] + cm[1][1]))
	fp := float64((cm[1][0] + cm[2][0]) + (cm[0][1] + cm[2][1]) + (cm[0][2] + cm[1][2]))
	fn := float64((cm[0][1] + cm[0][2]) + (cm[1][0] + cm[1][2]) + (cm[2][0] + cm[2][1]))
	// denominator :=
	return ((tp * tn) - (fp * fn)) / math.Sqrt((tp+fp)*(tp+fn)*(tn+fp)*(tn+fn))
}

func mcc2(cm [3][4]int) float64 {
	n := 2*(cm[0][0]*cm[1][1]-cm[1][0]*cm[0][1]) + 2*(cm[0][0]*cm[2][2]-cm[2][0]*cm[0][2]) + 2*(cm[1][1]*cm[2][2]-cm[2][1]*cm[1][2])
	d1 := (cm[0][0]+cm[1][0]+cm[2][0])*(cm[0][1]+cm[0][2]+cm[1][1]+cm[1][2]+cm[2][1]+cm[2][2]) + (cm[0][1]+cm[1][1]+cm[2][1])*(cm[0][0]+cm[0][2]+cm[1][0]+cm[1][2]+cm[2][0]+cm[2][2]) + (cm[0][2]+cm[1][2]+cm[2][2])*(cm[0][0]+cm[0][1]+cm[1][0]+cm[1][1]+cm[2][0]+cm[2][1])
	d2 := (cm[0][0]+cm[0][1]+cm[0][2])*(cm[1][0]+cm[1][1]+cm[1][2]+cm[2][0]+cm[2][1]+cm[2][2]) + (cm[1][0]+cm[1][1]+cm[1][2])*(cm[0][0]+cm[0][1]+cm[0][2]+cm[2][0]+cm[2][1]+cm[2][2]) + (cm[2][0]+cm[2][1]+cm[2][2])*(cm[0][0]+cm[0][1]+cm[0][2]+cm[2][0]+cm[2][1]+cm[2][2])

	return float64(n) / (math.Sqrt(float64(d1)) * math.Sqrt(float64(d2)))
}

func CBA(cm [3][4]int) float64 {
	nr := len(cm)
	np := len(cm[0])

	cba := make([]float64, nr)

	for i := 0; i < nr; i++ {
		ci_, c_i := 0.0, 0.0
		for j := 0; j < nr; j++ {
			ci_ += float64(cm[i][j])
			c_i += float64(cm[j][i])
		}
		if nr != np {
			ci_ += float64(cm[i][nr])
		}
		cba[i] = float64(cm[i][i]) / math.Max(ci_, c_i)
	}

	total := 0.0
	for t := 0; t < nr; t++ {
		total += cba[t]
	}
	return total / float64(nr)
}

// func score2(oc []int, end []rules.State, norm int) float64 {
// 	for i := 0; i < len(oc); i++ {
// 		if end[i] != rules.S_init && end[i] != rules.S__ {

// 		}
// 	}
// }

func step(previous, current, init, end *[]rules.State, cm *[3][4]int, occurrence *[]int, hocc, eocc *[]uint, rule *rules.Rule, use bool) {
	for c := 1; c < len(*init)-1; c++ {
		(*current)[c] = (*rule)[(*previous)[c-1]][(*previous)[c]][(*previous)[c+1]]
		if (*current)[c] == rules.S_init {
			(*current)[c] = (*init)[c]
		}
		if use {
			// ocurrence doesn't look to neighbors
			if (*current)[c] == (*end)[c] {
				(*occurrence)[c]++
			}
			// fmt.Println((*end)[c], rules.SS((*end)[c]), (*current)[c], rules.SS((*current)[c]))
			if rules.SS((*end)[c]) == "helix" {
				if rules.SS((*current)[c]) == "helix" {
					(*cm)[0][0]++
				} else if rules.SS((*current)[c]) == "strand" {
					(*cm)[0][1]++
				} else if rules.SS((*current)[c]) == "coil" {
					(*cm)[0][2]++
				} else {
					(*cm)[0][3]++
				}
			} else if rules.SS((*end)[c]) == "strand" {
				if rules.SS((*current)[c]) == "helix" {
					(*cm)[1][0]++
				} else if rules.SS((*current)[c]) == "strand" {
					(*cm)[1][1]++
				} else if rules.SS((*current)[c]) == "coil" {
					(*cm)[1][2]++
				} else {
					(*cm)[1][3]++
				}
			} else if rules.SS((*end)[c]) == "coil" {
				if rules.SS((*current)[c]) == "helix" {
					(*cm)[2][0]++
				} else if rules.SS((*current)[c]) == "strand" {
					(*cm)[2][1]++
				} else if rules.SS((*current)[c]) == "coil" {
					(*cm)[2][2]++
				} else {
					(*cm)[2][3]++
				}
			}
		}
	}

	// This occurrence look to neighbors *IF IN USE, COMMENT THE OCC CODE ABOVE
	// if use {
	// 	countOcc2(current, end, occurrence, hocc, eocc)
	// }
}

// func countOcc(curr, end *[]rules.State, occurrence *[]int) {
// 	var check bool
// 	for c := 1; c < len(*end)-1; c++ {
// 		switch (*end)[c] {
// 		case rules.S_e, rules.S_en, rules.S_ep, rules.S_eG, rules.S_eP, rules.S_eneg, rules.S_epos:
// 			if (c > 0) && (c < len(*end)-3) && !check {
// 				if testE3((*curr)[c], (*curr)[c+1], (*curr)[c+2]) {
// 					(*occurrence)[c]++
// 					break
// 				}
// 			}
// 			if (c > 1) && (c < len(*end)-2) && !check {
// 				if testE3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) {
// 					(*occurrence)[c]++
// 					break
// 				}
// 			}
// 			if (c > 2) && (c < len(*end)-1) && !check {
// 				if testE3((*curr)[c-2], (*curr)[c-1], (*curr)[c]) {
// 					(*occurrence)[c]++
// 					break
// 				}
// 			}
// 		case rules.S_h, rules.S_hn, rules.S_hp, rules.S_hG, rules.S_hP, rules.S_hneg, rules.S_hpos:
// 			if (c > 0) && (c < len(*end)-3) && !check {
// 				if testH3((*curr)[c], (*curr)[c+1], (*curr)[c+2]) {
// 					(*occurrence)[c]++
// 					break
// 				}
// 			}
// 			if (c > 1) && (c < len(*end)-2) && !check {
// 				if testH3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) {
// 					(*occurrence)[c]++
// 					break
// 				}
// 			}
// 			if (c > 2) && (c < len(*end)-1) && !check {
// 				if testH3((*curr)[c-2], (*curr)[c-1], (*curr)[c]) {
// 					(*occurrence)[c]++
// 					break
// 				}
// 			}
// 		case rules.S_c, rules.S_cn, rules.S_cp, rules.S_cG, rules.S_cP, rules.S_cneg, rules.S_cpos:
// 			if (c > 0) && (c < len(*end)-3) && !check {
// 				if testH3((*curr)[c], (*curr)[c+1], (*curr)[c+2]) || testE3((*curr)[c], (*curr)[c+1], (*curr)[c+2]) {
// 					break
// 				}
// 			}
// 			if (c > 1) && (c < len(*end)-2) && !check {
// 				if testH3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) || testE3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) {
// 					break
// 				}
// 			}
// 			if (c > 2) && (c < len(*end)-1) && !check {
// 				if testH3((*curr)[c-2], (*curr)[c-1], (*curr)[c]) || testE3((*curr)[c-2], (*curr)[c-1], (*curr)[c]) {
// 					break
// 				}
// 			}
// 			(*occurrence)[c]++
// 		}
// 	}
// }

// func countOcc2(curr, end *[]rules.State, occurrence *[]int, hocc *[]uint, eocc *[]uint) {
// 	// h := make([]uint, len(*curr))
// 	// e := make([]uint, len(*curr))
// 	for c := 1; c < len(*end)-1; c++ {
// 		if testE3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) {
// 			(*eocc)[c-1], (*eocc)[c], (*eocc)[c+1] = 1, 1, 1
// 		} else if testH3((*curr)[c-1], (*curr)[c], (*curr)[c+1]) {
// 			(*hocc)[c-1], (*hocc)[c], (*hocc)[c+1] = 1, 1, 1
// 		}
// 	}
// 	for c := 1; c < len(*end)-1; c++ {
// 		switch (*end)[c] {
// 		case rules.S_e, rules.S_en, rules.S_ep, rules.S_eG, rules.S_eP, rules.S_eneg, rules.S_epos:
// 			if (*eocc)[c] == 1 {
// 				(*occurrence)[c]++
// 			}
// 		case rules.S_h, rules.S_hn, rules.S_hp, rules.S_hG, rules.S_hP, rules.S_hneg, rules.S_hpos:
// 			if (*hocc)[c] == 1 {
// 				(*occurrence)[c]++
// 			}
// 		case rules.S_c, rules.S_cn, rules.S_cp, rules.S_cG, rules.S_cP, rules.S_cneg, rules.S_cpos:
// 			if (*eocc)[c] != 1 && (*hocc)[c] != 1 {
// 				(*occurrence)[c]++
// 			}
// 		}
// 		(*eocc)[c] = 0
// 		(*hocc)[c] = 0
// 	}
// }

// func testE3(p0, p1, p2 rules.State) bool {
// 	if p0 == rules.S_e || p0 == rules.S_en || p0 == rules.S_ep || p0 == rules.S_eG || p0 == rules.S_eP || p0 == rules.S_eneg || p0 == rules.S_epos {
// 		if p1 == rules.S_e || p1 == rules.S_en || p1 == rules.S_ep || p1 == rules.S_eG || p1 == rules.S_eP || p1 == rules.S_eneg || p1 == rules.S_epos {
// 			if p2 == rules.S_e || p2 == rules.S_en || p2 == rules.S_ep || p2 == rules.S_eG || p2 == rules.S_eP || p2 == rules.S_eneg || p2 == rules.S_epos {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// func testH3(p0, p1, p2 rules.State) bool {
// 	if p0 == rules.S_h || p0 == rules.S_hn || p0 == rules.S_hp || p0 == rules.S_hG || p0 == rules.S_hP || p0 == rules.S_hneg || p0 == rules.S_hpos {
// 		if p1 == rules.S_h || p1 == rules.S_hn || p1 == rules.S_hp || p1 == rules.S_hG || p1 == rules.S_hP || p1 == rules.S_hneg || p1 == rules.S_hpos {
// 			if p2 == rules.S_h || p2 == rules.S_hn || p2 == rules.S_hp || p2 == rules.S_hG || p2 == rules.S_hP || p2 == rules.S_hneg || p2 == rules.S_hpos {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }
