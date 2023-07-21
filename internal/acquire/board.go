package acquire

import "acquire/internal/util"

const BOARD_MAX_X = 12
const BOARD_MAX_Y = 9

type PlacedHotel struct {
	Hotel Hotel
	Pos   util.Point[int]
}

type Board struct {
	Matrix *util.Matrix[PlacedHotel]
}

func newBoard() *Board {
	return &Board{
		Matrix: util.NewMatrix[PlacedHotel](BOARD_MAX_X, BOARD_MAX_Y),
	}
}

func (b *Board) clone() *Board {
	return &Board{
		Matrix: b.Matrix.Copy(),
	}
}

func ValidPlacementPositions(game *Game) []util.Point[int] {
	positions := util.Map(game.CurrentPlayer().legalMoves(), func(val Tile) util.Point[int] {
		return val.Pos()
	})

	return positions
}
