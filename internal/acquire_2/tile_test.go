package acquire_2

import (
	"acquire/internal/util"
	"fmt"
	"testing"
)

func TestTilePos(t *testing.T) {
	testTile := func(tile Tile, x int, y int) {
		pos := tile.Pos()
		if pos != (util.Point[int]{X: x, Y: y}) {
			t.Fatal(fmt.Sprintf(
				"tile %s should be pos (%d, %d), it was (%d, %d)",
				tile.String(),
				x, y,
				pos.X, pos.Y,
			))
		}
	}

	testTile(Tile1A, 0, 0)
	testTile(Tile2A, 1, 0)
	testTile(Tile1I, 0, 8)

	testTile(Tile1B, 0, 1)
	testTile(Tile2B, 1, 1)
	testTile(Tile2I, 1, 8)

	testTile(Tile12A, 11, 0)
	testTile(Tile12B, 11, 1)
	testTile(Tile12I, 11, 8)
}

func TestBoardIndex(t *testing.T) {
	testTile := func(tile Tile, x int, y int) {
		idx := index(x, y)
		tileAtIdx := TileFromBoardIdx(idx)
		if tile != tileAtIdx {
			t.Fatal(fmt.Sprintf(
				"tile %s should be pos (%d, %d), it was (%d, %d)",
				tile.String(),
				x, y,
				Tile(idx).Pos().X, Tile(idx).Pos().Y,
			))
		}
	}

	testTile(Tile1A, 0, 0)
	testTile(Tile2A, 1, 0)
	testTile(Tile1I, 0, 8)

	testTile(Tile1B, 0, 1)
	testTile(Tile2B, 1, 1)
	testTile(Tile2I, 1, 8)

	testTile(Tile12A, 11, 0)
	testTile(Tile12B, 11, 1)
	testTile(Tile12I, 11, 8)
}
