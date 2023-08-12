package main

import (
	"acquire/internal/acquire"
	"acquire/internal/ai"
	"acquire/internal/util"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	runGame(0, 550, true)
}

func analyzeAIPerformance() {
	f, err := os.Create("./data_out.csv")
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)

	_ = w.Write([]string{"I", "T(s)", "P1 Int", "1", "2", "3", "4"})

	n := 100
	simCount := 1
	totalData := make([][]string, 0, n)
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			for s := 0; s < simCount; s++ {
				start := time.Now()
				intel := i*10 + 1
				game := runGame(s, intel, false)
				playerSlice := game.PlayerSlice()
				data := util.Map(playerSlice, func(player acquire.Player) string {
					return strconv.Itoa(player.NetWorth(game))
				})
				end := time.Now().Sub(start)
				data = append([]string{
					strconv.Itoa(i),
					strconv.Itoa(int(end.Seconds())),
					strconv.Itoa(intel)},
					data...,
				)
				totalData = append(totalData, data)
				_ = w.Write(data)
				fmt.Printf("I: %d, Int: %d, Sim #: %d, Game Turn: %d\n", i, intel, s, game.Turn)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	w.Flush()
}

func runGame(seed int, smartPlayerIntelligence int, display bool) *acquire.Game {
	rand.Seed(int64(seed))

	game := acquire.NewGame()

	agents := make(map[int]ai.IAgent)
	for _, player := range game.Players {
		agents[player.Id] = ai.NewStupidAgent()
	}
	//agents[game.Players[0].Id] = ai.NewHumanAgent()
	//agents[game.Players[0].Id] = ai.NewSmartAgent(smartPlayerIntelligence)
	//agents[game.Players[1].Id] = ai.NewSmartAgent(smartPlayerIntelligence)

	for !game.IsTerminal() {
		if display {
			acquire.Render(game)
		}

		agent := agents[game.CurrentPlayer().Id]
		actions := game.GetActions()
		action, err := agent.SelectAction(game, actions)
		if err != nil {
			panic(err)
		}

		if _action, ok := action.(acquire.IAction); ok {
			fmt.Println(_action.String(game))
		}

		newGame, _ := game.ApplyAction(action)
		game = newGame.(*acquire.Game)
	}

	if display {

		acquire.Render(game)

		fmt.Println()
		reason, end := game.CanEnd()
		game2 := game
		game2.Computed = acquire.NewComputed(game2)
		if !end {
			reason = "Game was forced to end."
		}
		fmt.Println("End Reason: " + reason)
	}

	return game
}
