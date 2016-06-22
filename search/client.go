package search

import (
	"fmt"
	"log"
	"math"
	"net/rpc"
	"time"

	"github.com/jgcarvalho/zeca-search/ca"
	"github.com/jgcarvalho/zeca-search/db"
)

func RunClient(conf Config) {
	client, err := rpc.DialHTTP("tcp", conf.Dist.MasterURL+":2222")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Le os dados das proteinas no DB
	fmt.Println("Loading proteins...")
	start, end, err := db.GetProteins(conf.DB)
	if err != nil {
		fmt.Println("Erro no banco de DADOS")
		panic(err)
	}

	var prob Probabilities
	// var rule rules.Rule
	var score float64
	var winner Individual
	var accepted, getnew bool
	accepted = true
	q := 1

	cellAuto := ca.Config{InitState: start, EndState: end, Steps: conf.CA.Steps, IgnoreSteps: conf.CA.IgnoreSteps}
	err = client.Call("MSG.GetProb", &q, &prob)
	if err != nil {
		log.Fatal("get prob error:", err)
	}

	g := 0
	for g < conf.EDA.Generations {
		err = client.Call("MSG.CheckProb", &prob.Generation, &getnew)
		if err != nil {
			log.Fatal("check prob error:", err)
		}
		if getnew {
			err = client.Call("MSG.GetProb", &q, &prob)
			if err != nil {
				log.Fatal("get prob error:", err)
			}
			// fmt.Println(prob.Data)
			// fmt.Println(prob.Generation)
		}

		if !accepted {
			if g != prob.Generation {
				g = prob.Generation
				accepted = true
			} else {
				// wait some time to get new prob
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}
		g = prob.Generation

		for i := 0; i < conf.EDA.Tournament; i++ {
			// fmt.Println("Generating rule")
			rule := GenRule(prob)
			// fmt.Println(prob.Data)
			// fmt.Println("Rule ok")
			// fmt.Println(start, end, tourn, ind, b, rule)
			score = cellAuto.Run(rule)
			fmt.Printf("Ind %d, score %f\n", i, score)
			if i == 0 {
				winner = Individual{Generation: prob.Generation, Rule: &rule, Score: score}
			} else {
				if score > winner.Score {
					winner = Individual{Generation: prob.Generation, Rule: &rule, Score: score}
				}
			}
		}

		// GOB DOES NOT ENCODE ZERO VALUES, SO ....
		if winner.Score == 0.0 {
			winner.Score = math.SmallestNonzeroFloat64
		}
		fmt.Println("winner", winner.Score)
		err = client.Call("MSG.SendWinner", winner, &accepted)
		if err != nil {
			log.Fatal("send winner error:", err)
		}
		// fmt.Println(prob.Generation, accepted)
	}
}
