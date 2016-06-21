package search

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"

	"github.com/gonum/stat"
)

type MSG int

type CurrentProb struct {
	sync.RWMutex
	Prob Probabilities
}

type Incoming struct {
	sync.Mutex
	N         int
	NMax      int
	Score     []float64
	NewProb   Probabilities
	BestScore float64
	Best      *Individual
}

var CurProb CurrentProb
var incoming Incoming

func RunServer(conf Config) {
	CurProb.Lock()
	CurProb.Prob.Generation = 0
	CurProb.Prob.Data = ReadProbRule(conf.Rules.Input)
	CurProb.Unlock()

	incoming.Lock()
	incoming.NMax = conf.EDA.Population / conf.EDA.Tournament
	incoming.Score = make([]float64, conf.EDA.Population/conf.EDA.Tournament)
	incoming.NewProb.Generation = 1
	incoming.NewProb.Data = ReadProbRule(conf.Rules.Input)
	incoming.NewProb.Data.Reset()
	incoming.BestScore = -math.MaxFloat64
	incoming.Unlock()

	msg := new(MSG)
	rpc.Register(msg)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":2222")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

}

func (t *MSG) CheckProb(gen *int, getnew *bool) error {
	CurProb.RLock()
	defer CurProb.RUnlock()
	if *gen != CurProb.Prob.Generation {
		*getnew = true
	} else {
		*getnew = false
	}
	return nil
}

func (t *MSG) GetProb(get *int, reply *Probabilities) error {
	CurProb.RLock()
	defer CurProb.RUnlock()
	*reply = Probabilities{Generation: CurProb.Prob.Generation, Data: CurProb.Prob.Data}
	// fmt.Println(*reply)
	// fmt.Println(reply)
	return nil
}

func (t *MSG) SendWinner(winner *Individual, accepted *bool) error {
	// testar se win gen e prob gen são iguais e testar se o numero de individuos não é maior que a população
	// if winner.Generation == CurrentProb.Generation &&
	incoming.Lock()
	defer incoming.Unlock()
	if incoming.N < incoming.NMax && winner.Generation == incoming.NewProb.Generation-1 {
		// fmt.Println(winner.Score)
		incoming.Score[incoming.N] = winner.Score
		incoming.NewProb.Data.Update(winner, incoming.NMax)
		if winner.Score > incoming.BestScore {
			incoming.BestScore = winner.Score
			incoming.Best = winner
		}
		incoming.N++
		*accepted = true
		if incoming.N == incoming.NMax {
			meanScore, stdScore := stat.MeanStdDev(incoming.Score, nil)
			fmt.Printf("G: %d, Mean Score: %.5f, StdDev Score: %.5f, Correct States: %.2f %%\n", winner.Generation, meanScore, stdScore, 100.0*math.Exp(meanScore))
			CurProb.Lock()
			CurProb.Prob.Save("prob_g" + strconv.Itoa(CurProb.Prob.Generation))
			CurProb.Prob.Generation = incoming.NewProb.Generation
			CurProb.Prob.Data.Copy(incoming.NewProb.Data)
			CurProb.Unlock()
			incoming.NewProb.Generation++
			incoming.NewProb.Data.Reset()
			incoming.N = 0
			fmt.Printf("BEST Score: %.5f, Correct States: %.2f %%\n", incoming.BestScore, 100.0*math.Exp(incoming.BestScore))
			// SaveBest(incoming.Best)

		}
	} else {
		*accepted = false
	}
	return nil
}
