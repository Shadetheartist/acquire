package internal

import (
	"acquire/internal/acquire_2"
	"acquire/internal/ai"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
)

func Benchmark(b *testing.B) {
	file, err := os.Create("../profile.prof")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	rand.Seed(int64(2))
	for i := 0; i < b.N; i++ {
		game := acquire_2.NewGame()

		agents := make(map[int]ai.IAgent)
		for _, player := range game.Players {
			agents[player.Id] = ai.NewStupidAgent()
		}

		for !game.IsTerminal() {

			agent := agents[game.CurrentPlayer().Id]
			actions := game.GetActions()
			action, err := agent.SelectAction(game, actions)
			if err != nil {
				panic(err)
			}
			newGame, _ := game.ApplyAction(action)
			game = newGame.(*acquire_2.Game)

		}
	}
}
