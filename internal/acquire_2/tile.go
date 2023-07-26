package acquire_2

import (
	"acquire/internal/util"
	"math/rand"
)

type Tile int

func (t Tile) String() string {
	return TileStringMap[t]
}

func TileFromString(str string) Tile {
	for idx, tileStr := range TileStringMap {
		if str == tileStr {
			return Tile(idx)
		}
	}

	return NoTile
}

func (t Tile) Index() int {
	// the tile IS the index of its position on the board if you
	// subtract one, since NoTile is the zero value
	return int(t - 1)
}

func (t Tile) Pos() util.Point[int] {
	idx := t.Index()
	return util.Point[int]{
		X: idx % BOARD_MAX_X,
		Y: idx / BOARD_MAX_X,
	}
}

func randomizedTiles() [108]Tile {
	tiles := [108]Tile{}

	for i := 0; i < len(TileList); i++ {
		tiles[i] = Tile(i)
	}

	rand.Shuffle(len(tiles), func(i, j int) {
		tiles[i], tiles[j] = tiles[j], tiles[i]
	})

	return tiles
}

const (
	NoTile Tile = iota
	Tile1A
	Tile2A
	Tile3A
	Tile4A
	Tile5A
	Tile6A
	Tile7A
	Tile8A
	Tile9A
	Tile10A
	Tile11A
	Tile12A
	Tile1B
	Tile2B
	Tile3B
	Tile4B
	Tile5B
	Tile6B
	Tile7B
	Tile8B
	Tile9B
	Tile10B
	Tile11B
	Tile12B
	Tile1C
	Tile2C
	Tile3C
	Tile4C
	Tile5C
	Tile6C
	Tile7C
	Tile8C
	Tile9C
	Tile10C
	Tile11C
	Tile12C
	Tile1D
	Tile2D
	Tile3D
	Tile4D
	Tile5D
	Tile6D
	Tile7D
	Tile8D
	Tile9D
	Tile10D
	Tile11D
	Tile12D
	Tile1E
	Tile2E
	Tile3E
	Tile4E
	Tile5E
	Tile6E
	Tile7E
	Tile8E
	Tile9E
	Tile10E
	Tile11E
	Tile12E
	Tile1F
	Tile2F
	Tile3F
	Tile4F
	Tile5F
	Tile6F
	Tile7F
	Tile8F
	Tile9F
	Tile10F
	Tile11F
	Tile12F
	Tile1G
	Tile2G
	Tile3G
	Tile4G
	Tile5G
	Tile6G
	Tile7G
	Tile8G
	Tile9G
	Tile10G
	Tile11G
	Tile12G
	Tile1H
	Tile2H
	Tile3H
	Tile4H
	Tile5H
	Tile6H
	Tile7H
	Tile8H
	Tile9H
	Tile10H
	Tile11H
	Tile12H
	Tile1I
	Tile2I
	Tile3I
	Tile4I
	Tile5I
	Tile6I
	Tile7I
	Tile8I
	Tile9I
	Tile10I
	Tile11I
	Tile12I
)

var TileList = []Tile{
	Tile1A,
	Tile2A,
	Tile3A,
	Tile4A,
	Tile5A,
	Tile6A,
	Tile7A,
	Tile8A,
	Tile9A,
	Tile10A,
	Tile11A,
	Tile12A,
	Tile1B,
	Tile2B,
	Tile3B,
	Tile4B,
	Tile5B,
	Tile6B,
	Tile7B,
	Tile8B,
	Tile9B,
	Tile10B,
	Tile11B,
	Tile12B,
	Tile1C,
	Tile2C,
	Tile3C,
	Tile4C,
	Tile5C,
	Tile6C,
	Tile7C,
	Tile8C,
	Tile9C,
	Tile10C,
	Tile11C,
	Tile12C,
	Tile1D,
	Tile2D,
	Tile3D,
	Tile4D,
	Tile5D,
	Tile6D,
	Tile7D,
	Tile8D,
	Tile9D,
	Tile10D,
	Tile11D,
	Tile12D,
	Tile1E,
	Tile2E,
	Tile3E,
	Tile4E,
	Tile5E,
	Tile6E,
	Tile7E,
	Tile8E,
	Tile9E,
	Tile10E,
	Tile11E,
	Tile12E,
	Tile1F,
	Tile2F,
	Tile3F,
	Tile4F,
	Tile5F,
	Tile6F,
	Tile7F,
	Tile8F,
	Tile9F,
	Tile10F,
	Tile11F,
	Tile12F,
	Tile1G,
	Tile2G,
	Tile3G,
	Tile4G,
	Tile5G,
	Tile6G,
	Tile7G,
	Tile8G,
	Tile9G,
	Tile10G,
	Tile11G,
	Tile12G,
	Tile1H,
	Tile2H,
	Tile3H,
	Tile4H,
	Tile5H,
	Tile6H,
	Tile7H,
	Tile8H,
	Tile9H,
	Tile10H,
	Tile11H,
	Tile12H,
	Tile1I,
	Tile2I,
	Tile3I,
	Tile4I,
	Tile5I,
	Tile6I,
	Tile7I,
	Tile8I,
	Tile9I,
	Tile10I,
	Tile11I,
	Tile12I,
}

var TileStringMap = []string{
	"None",
	"1A",
	"2A",
	"3A",
	"4A",
	"5A",
	"6A",
	"7A",
	"8A",
	"9A",
	"10A",
	"11A",
	"12A",
	"1B",
	"2B",
	"3B",
	"4B",
	"5B",
	"6B",
	"7B",
	"8B",
	"9B",
	"10B",
	"11B",
	"12B",
	"1C",
	"2C",
	"3C",
	"4C",
	"5C",
	"6C",
	"7C",
	"8C",
	"9C",
	"10C",
	"11C",
	"12C",
	"1D",
	"2D",
	"3D",
	"4D",
	"5D",
	"6D",
	"7D",
	"8D",
	"9D",
	"10D",
	"11D",
	"12D",
	"1E",
	"2E",
	"3E",
	"4E",
	"5E",
	"6E",
	"7E",
	"8E",
	"9E",
	"10E",
	"11E",
	"12E",
	"1F",
	"2F",
	"3F",
	"4F",
	"5F",
	"6F",
	"7F",
	"8F",
	"9F",
	"10F",
	"11F",
	"12F",
	"1G",
	"2G",
	"3G",
	"4G",
	"5G",
	"6G",
	"7G",
	"8G",
	"9G",
	"10G",
	"11G",
	"12G",
	"1H",
	"2H",
	"3H",
	"4H",
	"5H",
	"6H",
	"7H",
	"8H",
	"9H",
	"10H",
	"11H",
	"12H",
	"1I",
	"2I",
	"3I",
	"4I",
	"5I",
	"6I",
	"7I",
	"8I",
	"9I",
	"10I",
	"11I",
	"12I",
}
