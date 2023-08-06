package acquire

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

type Player struct {
	Id    int
	Money int
	Tiles [MAX_TILES_IN_HAND]Tile
	// mapped to hotels/stocks
	Stocks [NUM_CHAINS]int
}

func (player *Player) Name() string {
	return strconv.Itoa(player.Id)
}

// removeTileFromHand
// normally for use when the tile is placed on the board
// potentially could be used for removing permanently unplayable tiles from the game
func (player *Player) removeTileFromHand(tile Tile) error {

	// put the tile in the first empty slot
	for idx, t := range player.Tiles {
		if t == tile {
			player.Tiles[idx] = NoTile
			return nil
		}
	}

	return errors.New("cannot remove tile from hand, tile is not in hand")
}

func (player *Player) takeTileFromBank(game *Game) error {

	// look through the bank tiles array until a valid tile is found, that is the tile the player will take
	var tile Tile
	var bankIdx int
	for bankIdx, tile = range game.Tiles {
		if tile != NoTile {
			break
		}
	}

	if tile == NoTile {
		return errors.New("player cannot take a tile from the bank, the bank has no tiles remaining")
	}

	// put the tile in the first empty slot
	for idx, t := range player.Tiles {
		if t == NoTile {
			player.Tiles[idx] = tile

			// only set the bank
			game.Tiles[bankIdx] = NoTile
			return nil
		}
	}

	return errors.New("player cannot take a tile from the bank, their hand is full")
}

func (player *Player) returnTileToBank(game *Game, tile Tile) error {

	// look through the bank tiles array until a NoTile slot is found,
	// this will be where we replace the tile in the array
	var bankIdx int
	var bankTile Tile
	for bankIdx, bankTile = range game.Tiles {
		if bankTile == NoTile {
			break
		}
	}

	if bankTile != NoTile {
		return errors.New("player cannot replace a tile into the bank, the bank is full somehow")
	}

	err := player.removeTileFromHand(tile)
	if err != nil {
		return err
	}

	game.Tiles[bankIdx] = tile

	return nil
}

func (player *Player) takeStockFromBank(game *Game, chain Hotel, amount int) error {
	idx := chain.Index()

	if game.Stocks[idx] < amount {
		return errors.New("can't take stock from bank, not enough stock in bank to take")
	}

	game.Stocks[idx] -= amount
	player.Stocks[idx] += amount

	return nil
}

func (player *Player) giveStockToBank(game *Game, chain Hotel, amount int) error {
	idx := chain.Index()

	if player.Stocks[idx] < amount {
		return errors.New("can't give stock to bank, player does not have the requested amount to give")
	}

	game.Stocks[idx] += amount
	player.Stocks[idx] -= amount

	return nil
}

func (player *Player) returnedAmount(tradeInAmount int) int {
	return tradeInAmount / 2
}

func (player *Player) canTradeIn(game *Game, in Hotel, out Hotel, tradeInAmount int) error {

	if (tradeInAmount % 2) != 0 {
		return errors.New("trade in amount must be a multiple of two")
	}

	if player.Stocks[in] < tradeInAmount {
		return errors.New("player does not have enough stock to trade in for")
	}

	returnAmount := player.returnedAmount(tradeInAmount)
	if returnAmount < 1 {
		return nil
	}

	outRemaining := game.Stocks[out.Index()]

	if outRemaining < returnAmount {
		return errors.New("there is not enough stock to trade-in for")
	}

	return nil
}

func (player *Player) tradeIn(game *Game, in Hotel, out Hotel, tradeInAmount int) error {

	err := player.canTradeIn(game, in, out, tradeInAmount)
	if err != nil {
		return err
	}

	returnAmount := player.returnedAmount(tradeInAmount)

	// return the stock to the bank
	_ = player.giveStockToBank(game, in, tradeInAmount)

	// take half the amount of the new stock
	_ = player.takeStockFromBank(game, out, returnAmount)

	return nil
}

func (player *Player) canSellStock(stock Stock, amount int) error {
	if player.Stocks[Hotel(stock).Index()] < amount {
		return errors.New("amount sold cannot be greater than the amount the player has")
	}

	return nil
}

func (player *Player) sellStock(game *Game, stock Stock, amount int) error {
	err := player.canSellStock(stock, amount)
	if err != nil {
		return err
	}

	hotel := Hotel(stock)
	err = player.giveStockToBank(game, hotel, amount)
	if err != nil {
		panic("shouldn't give more than there are in player's inventory")
	}

	chainSize := game.ChainSize[hotel.Index()]
	value := sharesCalc(hotel, chainSize, amount)

	player.Money += value

	return nil
}

func (player *Player) buyStock(game *Game, hotel Hotel, amount int) error {

	chainSize := game.ChainSize[hotel.Index()]

	if chainSize == 0 {
		return errors.New("chain does not exist at the moment")
	}

	cost := sharesCalc(hotel, chainSize, amount)

	err := player.pay(cost)
	if err != nil {
		return err
	}

	amountAvailable := game.Stocks[hotel.Index()]
	if amountAvailable < amount {
		return fmt.Errorf("can't buy %d stocks in %s, there's only %d remaining", amount, hotel.String(), amountAvailable)
	}

	err = player.takeStockFromBank(game, hotel, amount)
	if err != nil {
		panic(err)
	}

	return nil
}

func (player *Player) pay(amount int) error {
	if player.Money < amount {
		return fmt.Errorf("player '%d' cannot afford to pay $%d", player.Id, amount)
	}

	player.Money += amount

	return nil
}

// refreshTiles
// when a player has no legal moves left to play, they can refresh their hand with this func
// puts all tiles back in the Game inv, then takes 6 new ones
func (player *Player) refreshTiles(game *Game) {

	for _, t := range player.Tiles {
		err := player.returnTileToBank(game, t)
		if err != nil {
			panic(err)
		}
	}

	// shuffle the tiles
	rand.Shuffle(len(game.Tiles), func(i, j int) {
		game.Tiles[i], game.Tiles[j] = game.Tiles[j], game.Tiles[i]
	})

	for i := 0; i < 6; i++ {
		err := player.takeTileFromBank(game)
		if err != nil {
			game.end("no tiles left")
			return
		}
	}
}

func (player *Player) NetWorth(game *Game) int {
	return game.Computed.PlayerNetWorth[player.Id]
}
