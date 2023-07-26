package console_interface

import (
	"acquire/internal/acquire_2"
	"acquire/internal/util"
	"fmt"
	"strconv"
	"strings"
)

var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}

const FILL_SIZE = 3

func Render(game *acquire_2.Game) {
	fmt.Print("\033[H\033[2J")
	fmt.Println()
	fmt.Println("Acquire | Tiles Left:", game.NumRemainingTiles())
	renderPlayers(game)
	renderPlayerInventories(game)
	renderBoard(game)
	renderCurrentPlayer(game)
	println()
}

func fill(s string, n int) string {
	l := n - len([]rune(s))
	if l < 0 {
		l = 0
	}

	return s + strings.Repeat(" ", l)
}
func rfill(s string, n int) string {
	l := n - len([]rune(s))
	if l < 0 {
		l = 0
	}

	return strings.Repeat(" ", l) + s
}

func renderCurrentPlayer(game *acquire_2.Game) {
	player := game.CurrentPlayer()

	fmt.Println()

	fmt.Printf(fill("Tiles:", 8))
	for _, t := range player.Tiles {
		print(fill(fmt.Sprintf("%s ", t.String()), 4))
	}
	fmt.Println()

}

func renderPlayers(game *acquire_2.Game) {
	fmt.Printf(fill("Turn: ", 8))
	fillSize := 10
	for _, p := range game.Players {
		if p.Id == game.CurrentPlayer().Id {
			print(fill(fmt.Sprintf("[%s] ", p.Name()), fillSize))
			continue
		}
		print(fill(fmt.Sprintf("%s ", p.Name()), fillSize))
	}
	fmt.Println()
}

func renderPlayerInventories(game *acquire_2.Game) {
	fmt.Println("Inventories:")
	nameFillSize := 9
	fillSize := 3

	print(fill(" ", nameFillSize))
	print(rfill("", nameFillSize))
	for _, h := range acquire_2.HotelChainList {
		print(fill(h.Initial(), fillSize))
	}
	println()

	for _, p := range game.Players {
		print(rfill(p.Name()+" ", nameFillSize))
		print(fill(fmt.Sprintf("$%d", p.Money), nameFillSize))
		for h, _ := range acquire_2.HotelChainList {
			print(fill(strconv.Itoa(p.Stocks[h]), fillSize))
		}
		println()
	}
	fmt.Println()
}

func renderBoard(game *acquire_2.Game) {

	for x := 0; x <= acquire_2.BOARD_MAX_X; x++ {
		if x == 0 {
			print(fill(" ", FILL_SIZE))
			continue
		}
		print(fill(strconv.Itoa(x), FILL_SIZE))
	}
	println()

	validPlacementPositions := util.Map(game.Computed.LegalMoves, func(val acquire_2.Tile) util.Point[int] {
		return val.Pos()
	})
	isValidPlacement := func(x int, y int) bool {
		_, ok := util.IndexOf(validPlacementPositions, util.Point[int]{X: x, Y: y})
		return ok
	}

	for y := 0; y < acquire_2.BOARD_MAX_Y; y++ {
		print(chars[y] + strings.Repeat(" ", FILL_SIZE-1))

		for x := 0; x < acquire_2.BOARD_MAX_X; x++ {
			placedHotel := game.PlacementAtPos(util.Point[int]{X: x, Y: y})

			if placedHotel.Hotel == acquire_2.NoHotel {
				if isValidPlacement(x, y) {
					print(fill("\u25CB", FILL_SIZE))
				} else {
					print(fill("\u25A1", FILL_SIZE))
				}
				continue
			}

			if placedHotel.Hotel == acquire_2.UndefinedHotel {
				print(fill("\u25A0", FILL_SIZE))
				continue
			}

			print(fill(placedHotel.Hotel.Initial(), FILL_SIZE))
		}
		println()
	}
}
