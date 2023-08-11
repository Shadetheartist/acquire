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

		// replace tile on board
		if currentPlacedHotel.Hotel != rootHotel.Hotel {
			game.Board[currentPlacedHotel.Tile.Index()] = PlacedHotel{
				Hotel: rootHotel.Hotel,
				Tile:  currentPlacedHotel.Tile,
			}

			// chain got bigger
			if currentPlacedHotel.Hotel != UndefinedHotel {
				game.modifyChainSize(currentPlacedHotel.Hotel, -1)
			}

			game.modifyChainSize(rootHotel.Hotel, 1)
		}

		// add neighbors to stack
		neighborPts := currentPlacedHotel.Tile.Pos().OrthogonalNeighbours()

		for _, pt := range neighborPts {

			if !isInBounds(pt.X, pt.Y) {
				continue
			}

			neighbourIdx := index(pt.X, pt.Y)
			neighbor := game.Board[neighbourIdx]

			// don't add NoHotel neighbors to the stack
			if neighbor.Hotel == NoHotel {
				continue
			}

			//avoid revisiting
			if visited[neighbourIdx] {
				continue
			}

			// the hotel we're already propagating, no need to propagate this direction
			if neighbor.Hotel == rootHotel.Hotel {
				continue
			}

			visited[neighbourIdx] = true
			stack = append(stack, neighbor)
		}
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
