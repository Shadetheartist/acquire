package acquire

import (
	"acquire/internal/util"
)

// propagateHotelChain
// this breadth-first searches the matrix for connected hotels of
// any chain and converts them to the type of the rootHotel
func propagateHotelChain(game *Game, rootHotel PlacedHotel) {
	visited := make([]bool, len(game.Board))
	stack := make([]PlacedHotel, 0, 8)

	// root stack push
	stack = append(stack, rootHotel)
	var currentPlacedHotel PlacedHotel

	for len(stack) > 0 {
		// pop off stack
		currentPlacedHotel, stack = stack[len(stack)-1], stack[:len(stack)-1]
		idx := int(currentPlacedHotel.Tile)

		if currentPlacedHotel.Hotel == NoHotel {
			continue
		}

		//avoid revisiting
		if visited[idx] {
			continue
		}

		//keep track of what has been visited already
		visited[idx] = true

		// the hotel we're already propagating
		if currentPlacedHotel.Hotel == rootHotel.Hotel && currentPlacedHotel != rootHotel {
			continue
		}

		// replace tile on board
		game.Board[currentPlacedHotel.Tile.Index()] = PlacedHotel{
			Hotel: rootHotel.Hotel,
			Tile:  currentPlacedHotel.Tile,
		}
		// track chain size
		game.ChainSize[rootHotel.Hotel.Index()]++
		if currentPlacedHotel.Hotel != UndefinedHotel {
			game.ChainSize[currentPlacedHotel.Hotel.Index()]--
		}

		// add neighbors to stack
		neighborPts := currentPlacedHotel.Tile.Pos().OrthogonalNeighbours()
		neighbors := util.Map(neighborPts, func(val util.Point[int]) PlacedHotel {
			if !isInBounds(val.X, val.Y) {
				return PlacedHotel{}
			}

			idx := index(val.X, val.Y)

			neighbor := game.Board[idx]

			return neighbor
		})

		stack = append(stack, neighbors...)

	}
}

func isInBounds(x int, y int) bool {
	return x >= 0 && x < BOARD_MAX_X && y >= 0 && y < BOARD_MAX_Y
}

func index(x int, y int) int {
	return (y * BOARD_MAX_X) + x
}

func getNeighbors(game *Game, pt util.Point[int]) []PlacedHotel {
	neighbors := make([]PlacedHotel, len(util.Directions))

	for d, npt := range pt.OrthogonalNeighbours() {
		if !isInBounds(npt.X, npt.Y) {
			neighbors[d] = PlacedHotel{}
			continue
		}

		idx := index(npt.X, npt.Y)
		nb := game.Board[idx]
		neighbors[d] = nb
	}

	return neighbors
}
