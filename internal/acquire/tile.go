package acquire

import (
	"acquire/internal/util"
	"fmt"
	"math/rand"
	"strconv"
	"unicode"
)

type Tile string

func (t Tile) String() string {
	return string(t)
}

func (t Tile) Pos() util.Point[int] {
	x, y, err := parseTileIntoCoords(string(t))
	if err != nil {
		panic(err)
	}
	return util.Point[int]{
		X: x - 1,
		Y: y - 1,
	}
}

func randomizedTiles() []Tile {
	tiles := make([]Tile, len(TileList))

	copy(tiles, TileList)

	rand.Shuffle(len(tiles), func(i, j int) {
		tiles[i], tiles[j] = tiles[j], tiles[i]
	})

	return tiles
}

// isActualTile
// this compares the tile (which is just a string) to the list
// of predefined tiles to see if the tile is a valid board space
func isActualTile(tile Tile) bool {
	for _, t := range TileList {
		if t == tile {
			return true
		}
	}

	return false
}

func isLegalToPlace(game *Game, tile Tile) bool {
	pos := tile.Pos()
	neighboringHotels := game.Board.Matrix.GetNeighbors(pos)
	chainsInNeighbors := getChainsInNeighbors(neighboringHotels)

	// this tile would start a merger if placed
	if len(chainsInNeighbors) > 1 {
		// if both neighbors are size > 10, then they are both safe from merging and the placement isn't legal
		sizes := countHotelChains(game, chainsInNeighbors)
		anyMergable := util.AnyInMap(sizes, func(key Hotel, val int) bool {
			return val < 11
		})

		if !anyMergable {
			return false
		}
	}

	// this would found a new chain if placed
	undefinedNeighbors := getUndefinedNeighbors(neighboringHotels)
	if len(undefinedNeighbors) > 0 {
		// if there are no available chains left to create, this move is invalid
		availableChains := getAvailableHotelChains(game)
		if len(availableChains) < 1 {
			return false
		}
	}

	return true
}

func parseTileIntoCoords(input string) (int, int, error) {
	numStr := ""
	letter := ' '

	for _, char := range input {
		if unicode.IsDigit(char) {
			numStr += string(char)
		} else if unicode.IsLetter(char) {
			letter = unicode.ToUpper(char)
		}
	}

	if numStr == "" || letter == ' ' {
		return 0, 0, fmt.Errorf("invalid input")
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, 0, err
	}

	// 'A' corresponds to 1, 'B' corresponds to 2, etc.
	letterValue := int(letter - 'A' + 1)

	return num, letterValue, nil
}

const (
	NoTile  Tile = ""
	Tile1A  Tile = "1A"
	Tile2A  Tile = "2A"
	Tile3A  Tile = "3A"
	Tile4A  Tile = "4A"
	Tile5A  Tile = "5A"
	Tile6A  Tile = "6A"
	Tile7A  Tile = "7A"
	Tile8A  Tile = "8A"
	Tile9A  Tile = "9A"
	Tile10A Tile = "10A"
	Tile11A Tile = "11A"
	Tile12A Tile = "12A"
	Tile1B  Tile = "1B"
	Tile2B  Tile = "2B"
	Tile3B  Tile = "3B"
	Tile4B  Tile = "4B"
	Tile5B  Tile = "5B"
	Tile6B  Tile = "6B"
	Tile7B  Tile = "7B"
	Tile8B  Tile = "8B"
	Tile9B  Tile = "9B"
	Tile10B Tile = "10B"
	Tile11B Tile = "11B"
	Tile12B Tile = "12B"
	Tile1C  Tile = "1C"
	Tile2C  Tile = "2C"
	Tile3C  Tile = "3C"
	Tile4C  Tile = "4C"
	Tile5C  Tile = "5C"
	Tile6C  Tile = "6C"
	Tile7C  Tile = "7C"
	Tile8C  Tile = "8C"
	Tile9C  Tile = "9C"
	Tile10C Tile = "10C"
	Tile11C Tile = "11C"
	Tile12C Tile = "12C"
	Tile1D  Tile = "1D"
	Tile2D  Tile = "2D"
	Tile3D  Tile = "3D"
	Tile4D  Tile = "4D"
	Tile5D  Tile = "5D"
	Tile6D  Tile = "6D"
	Tile7D  Tile = "7D"
	Tile8D  Tile = "8D"
	Tile9D  Tile = "9D"
	Tile10D Tile = "10D"
	Tile11D Tile = "11D"
	Tile12D Tile = "12D"
	Tile1E  Tile = "1E"
	Tile2E  Tile = "2E"
	Tile3E  Tile = "3E"
	Tile4E  Tile = "4E"
	Tile5E  Tile = "5E"
	Tile6E  Tile = "6E"
	Tile7E  Tile = "7E"
	Tile8E  Tile = "8E"
	Tile9E  Tile = "9E"
	Tile10E Tile = "10E"
	Tile11E Tile = "11E"
	Tile12E Tile = "12E"
	Tile1F  Tile = "1F"
	Tile2F  Tile = "2F"
	Tile3F  Tile = "3F"
	Tile4F  Tile = "4F"
	Tile5F  Tile = "5F"
	Tile6F  Tile = "6F"
	Tile7F  Tile = "7F"
	Tile8F  Tile = "8F"
	Tile9F  Tile = "9F"
	Tile10F Tile = "10F"
	Tile11F Tile = "11F"
	Tile12F Tile = "12F"
	Tile1G  Tile = "1G"
	Tile2G  Tile = "2G"
	Tile3G  Tile = "3G"
	Tile4G  Tile = "4G"
	Tile5G  Tile = "5G"
	Tile6G  Tile = "6G"
	Tile7G  Tile = "7G"
	Tile8G  Tile = "8G"
	Tile9G  Tile = "9G"
	Tile10G Tile = "10G"
	Tile11G Tile = "11G"
	Tile12G Tile = "12G"
	Tile1H  Tile = "1H"
	Tile2H  Tile = "2H"
	Tile3H  Tile = "3H"
	Tile4H  Tile = "4H"
	Tile5H  Tile = "5H"
	Tile6H  Tile = "6H"
	Tile7H  Tile = "7H"
	Tile8H  Tile = "8H"
	Tile9H  Tile = "9H"
	Tile10H Tile = "10H"
	Tile11H Tile = "11H"
	Tile12H Tile = "12H"
	Tile1I  Tile = "1I"
	Tile2I  Tile = "2I"
	Tile3I  Tile = "3I"
	Tile4I  Tile = "4I"
	Tile5I  Tile = "5I"
	Tile6I  Tile = "6I"
	Tile7I  Tile = "7I"
	Tile8I  Tile = "8I"
	Tile9I  Tile = "9I"
	Tile10I Tile = "10I"
	Tile11I Tile = "11I"
	Tile12I Tile = "12I"
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
