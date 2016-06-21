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
			step(&previous, &current, &init, &end, &occurrence, &rule, use)
			// fmt.Println(current)
		} else {
			step(&current, &previous, &init, &end, &occurrence, &rule, use)
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
	return sc / float64(valid)
}

func step(previous, current, init, end *[]rules.State, occurrence *[]int, rule *rules.Rule, use bool) {
	for c := 1; c < len(*init)-1; c++ {
		(*current)[c] = (*rule)[(*previous)[c-1]][(*previous)[c]][(*previous)[c+1]]
		if (*current)[c] == rules.S_init {
			(*current)[c] = (*init)[c]
		}
		if use {
			// need change to consider neighbors
			if (*current)[c] == (*end)[c] || (*end)[c] == rules.S_init {
				(*occurrence)[c]++
			}
		}
		// (*current)[c] = (*rule)[rules.Pattern{(*previous)[c-1], (*previous)[c], (*previous)[c+1]}]
		// if string((*current)[c][0]) == "?" {
		// 	(*current)[c] = (*init)[c]
		// }
		// if use {
		// 	if (*current)[c] == (*end)[c] || string((*end)[c]) == "?" {
		// 		(*occurrence)[c]++
		// 	}
		// }
	}
}
