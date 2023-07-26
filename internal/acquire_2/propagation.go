package acquire_2

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
	var placedHotel PlacedHotel

	for len(stack) > 0 {
		// pop off stack
		placedHotel, stack = stack[len(stack)-1], stack[:len(stack)-1]
		idx := int(placedHotel.Tile)

		if placedHotel.Hotel == NoHotel {
			continue
		}

		//avoid revisiting
		if visited[idx] {
			continue
		}

		//keep track of what has been visited already
		visited[idx] = true

		// the hotel we're already propagating
		if placedHotel.Hotel == rootHotel.Hotel && placedHotel != rootHotel {
			continue
		}

		game.Board[placedHotel.Tile.Index()] = PlacedHotel{
			Hotel: rootHotel.Hotel,
			Tile:  placedHotel.Tile,
		}

		// add neighbors to stack
		neighborPts := placedHotel.Tile.Pos().OrthogonalNeighbours()
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
