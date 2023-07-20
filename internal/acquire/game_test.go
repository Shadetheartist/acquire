package acquire

import (
	"log"
	"os"
	"runtime/pprof"
	"testing"
)

// profiling reveals that most of the cpu tile is spent recalculating the network of chained hotels
// which was totally expected lol

func Benchmark(b *testing.B) {

	file, err := os.Create("../../profile.prof")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	inputInterface := &MockInputInterface{}

	for i := 0; i < b.N; i++ {

		game := NewGame(inputInterface)

		for {
			game.Step()
			if game.IsOver {
				break
			}
		}
	}
}
