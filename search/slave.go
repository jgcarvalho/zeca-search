package search

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jgcarvalho/zeca-search/ca"
	"github.com/jgcarvalho/zeca-search/db"
	zmq "github.com/pebbe/zmq4"
)

func RunSlave(conf Config) {

	// Cria o receptor que recebe a probabilidade emitida pelo master na porta A
	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Connect("tcp://" + conf.Dist.MasterURL + ":" + conf.Dist.PortA)

	// Cria o emissor que envia o individuo vencedor do torneio na rede pela
	// porta B
	sender, _ := zmq.NewSocket(zmq.PUSH)
	defer sender.Close()
	sender.Connect("tcp://" + conf.Dist.MasterURL + ":" + conf.Dist.PortB)

	// semente randomica
	rand.Seed(time.Now().UTC().UnixNano())

	// Le os dados das proteinas no DB
	fmt.Println("Loading proteins...")
	start, end, err := db.GetProteins(conf.DB)
	if err != nil {
		fmt.Println("Erro no banco de DADOS")
		panic(err)
	}

	// ALERT THE CODE BELOW IS ONLY FOR TEST
	// start, end := []string{"#", "M", "A", "D", "F", "G", "H", "I", "K", "#", "A", "A", "#"},
	// 	[]string{"#", "_", "_", "*", "*", "*", "*", "_", "|", "#", "|", "_", "#"}
	// start, end := []string{"#", "In", "An", "Dp", "Fn", "GG", "Ln", "In", "Kp", "#", "Ep", "PP", "#"},
	// 	[]string{"#", "_n", "_n", "*p", "*n", "*G", "*n", "_n", "|p", "#", "|p", "_P", "#"}
	// start, end := []string{"#", "Ep", "PP", "#"},
	// 	[]string{"#", "|p", "_P", "#"}

	fmt.Println("Done")

	var prob Probabilities

	// var tourn Tournament
	// tourn = make([]Individual, conf.EDA.Tournament)

	var (
		winner Individual
		// b      []byte
		m      []byte
		conerr error
	)

	cellAuto := ca.Config{InitState: start, EndState: end, Steps: conf.CA.Steps, IgnoreSteps: conf.CA.IgnoreSteps}

	var read bytes.Buffer
	// read.Grow(1000000)
	gob.RegisterName("probs", Probabilities{})
	decoder := gob.NewDecoder(&read)

	var write bytes.Buffer
	gob.RegisterName("winner", Individual{})
	encoder := gob.NewEncoder(&write)

	for {
		// m é a mensagem com as probabilidades
		m, conerr = receiver.RecvBytes(0)
		if conerr == nil {
			// read := bytes.NewBuffer(m)
			// decoder := gob.NewDecoder(read)
			read.Reset()
			read.Write(m)
			decoder.Decode(&prob)

			// json.Unmarshal([]byte(m), &prob)
			// fmt.Printf("PID: %d, Geracacao: %d\n", prob.PID, prob.Generation)

			for i := 0; i < conf.EDA.Tournament; i++ {
				rule := GenRule(prob)
				// fmt.Println(rule)
				// fmt.Println(start, end, tourn, ind, b, rule)
				score := cellAuto.Run(rule)
				// fmt.Println(score)
				if i == 0 {
					winner = Individual{PID: prob.PID, Generation: prob.Generation, Rule: &rule, Score: score}
				} else {
					if score > winner.Score {
						winner = Individual{PID: prob.PID, Generation: prob.Generation, Rule: &rule, Score: score}
					}
				}
			}

			// GOB DOES NOT ENCODE ZERO VALUES, SO ....
			if winner.Score == 0.0 {
				winner.Score = math.SmallestNonzeroFloat64
			}

			// write := new(bytes.Buffer)
			// encoder := gob.NewEncoder(write)
			write.Reset()
			encoder.Encode(winner)
			sender.SendBytes(write.Bytes(), 0)
			// b, _ = json.Marshal(winner)
			fmt.Println("Winner:", winner.Score)
			// sender.Send(string(b), 0)
		} else {
			// Erro na conexão
			fmt.Println(conerr)
		}
	}

	// para cada probabilidade recebida
	//criar t individuos do torneio
	// contruir e rodar automato celular
	// computar o Score

}
