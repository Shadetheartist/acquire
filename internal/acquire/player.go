package acquire

import (
	"acquire/internal/util"
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

/* player turn

1: play a tile onto the board
2: buy up to 3 stocks
3: draw a new random tile

--

when a tile is placed it can have one of these effects

- found a new hotel chain
- grow a hotel chain
- merge 2 or more hotel chains
- no effect (first tile, no connections)

--

once a tile is placed next to another tile that isn't already part of a hotel chain,
a new chain is created (player's choice)

you get 1 free stock in the chain for founding it ( if available )

if there are no more hotels left to found, you cannot play a tile to which would found a new one


--

when you place a tile next to an existing hotel chain, that hotel chain grows.

if the new tile is touching two or more hotel chains, a 'merger' takes place. The bigger hotel chain acquires the smaller one.
In a tie of hotel size, the player decides which hotel is acquired
During the merger, the tile you place does not count toward the size of either hotel.

As part of the merger, each player counts their stocks to determine the majority and minority shareholders
if tied, add the bonuses together and split them evenly, rounded to the nearest 100


*/

type Player struct {
	PlayerName string
	agent      IAgent
	game       *Game
	Inventory  *Inventory
}

func (p *Player) Name() string {
	return p.PlayerName
}

func NewPlayer(game *Game, name string, agentFactoryFn func(player *Player) IAgent) *Player {
	player := &Player{
		PlayerName: name,
		game:       game,
	}

	player.agent = agentFactoryFn(player)

	player.Inventory = newInventory(player, 0)

	return player
}

func (p *Player) pay(amount int) error {
	if p.Inventory.Money < amount {
		return fmt.Errorf("player '%s' cannot afford to pay $%d", p.PlayerName, amount)
	}

	p.game.Inventory.takeMoney(p.Inventory, amount)

	return nil
}

func (p *Player) buyStock(hotel Hotel, amount int) error {

	if !isActualHotelChain(hotel) {
		panic("trying to buy a non-chain hotel")
	}

	chainSize := countHotelChain(p.game, hotel)

	if chainSize == 0 {
		return errors.New("chain does not exist yet")
	}

	cost := sharesCalc(hotel, chainSize, amount)

	err := p.pay(cost)
	if err != nil {
		return err
	}

	amountAvailable := len(p.game.Inventory.Stocks[hotel].Items)
	if amountAvailable < amount {
		return fmt.Errorf("can't buy %d stocks in %s, there's only %d remaining", amount, hotel.String(), amountAvailable)
	}

	amountTaken := 0
	for i := 0; i < amount; i++ {
		err := p.Inventory.Stocks[hotel].take(p.game.Inventory.Stocks[hotel])
		if err != nil {
			return err
		}
		amountTaken++
	}

	return nil
}

func (p *Player) placeTile(tile Tile) {

	// move the tile from the player's inventory to the board
	p.Inventory.Tiles.remove(tile)
	pos := tile.Pos()
	matrix := p.game.Board.Matrix
	matrix.Set(pos.X, pos.Y, PlacedHotel{Pos: pos, Hotel: UndefinedHotel})

	neighboringHotels := matrix.GetNeighbors(pos)

	// no effect occurs when there are no adjacent tiles to the placed tile
	// so check if that's the case and early exit

	// no neighbors - no effect
	if !hasNeighboringHotel(neighboringHotels) {
		return
	}

	// growing a chain - occurs when, of all neighbors, there is only one type of hotel
	chainsInNeighbors := getChainsInNeighbors(neighboringHotels)

	if len(chainsInNeighbors) == 1 {
		hotel := chainsInNeighbors[0]
		newPlacedHotel := PlacedHotel{Pos: pos, Hotel: hotel}
		matrix.Set(pos.X, pos.Y, newPlacedHotel)

		propagateHotelChain(p.game, newPlacedHotel)

		// only one effect can occur per placement
		return
	}

	// merger - if there are more than two chains in the neighboring tiles, a merger must take place
	if len(chainsInNeighbors) > 1 {

		largestChains, _ := getLargestChainsOf(p.game, chainsInNeighbors)

		acquiringHotel := largestChains[0]
		var err error
		if len(largestChains) > 1 {
			acquiringHotel, err = p.agent.DetermineHotelToMerge(largestChains)
			if err != nil {
				panic(err)
			}
		}

		// count the hotel chains *before* the tile is placed for accurate stock purchasing
		chainSizeMap := countHotelChains(p.game, chainsInNeighbors)

		// in order of largest to smallest
		sorted := sortChainSizeMap(chainSizeMap)

		// pay-out shareholders, replace hotel chain (implicit, as this is always calculated at runtime)
		for _, chain := range sorted {

			acquiredHotel := chain.Hotel

			// don't sell the take-over hotel lol
			if acquiredHotel == acquiringHotel {
				continue
			}

			p.game.payShareholderBonuses(acquiredHotel)

			// players, in order starting with the current player,
			// decide to hold, trade, or sell stocks in the acquired chains
			players := p.game.playersInTurnOrder()
			for _, p := range players {
				action, err := p.agent.DetermineMergerAction(acquiredHotel)
				if err != nil {
					panic(err)
				}

				switch action {
				case Hold:
					break
				case Trade:
					amount, err := p.agent.DetermineTradeInAmount(acquiredHotel, acquiringHotel)
					if err != nil {
						panic(err)
					}

					err = p.tradeIn(Stock(acquiredHotel), Stock(acquiringHotel), amount)
					if err != nil {
						panic(err)
					}
					break

				case Sell:
					amount, err := p.agent.DetermineStockSellAmount(acquiredHotel)
					if err != nil {
						panic(err)
					}

					err = p.sellStock(Stock(acquiredHotel), amount)
					if err != nil {
						panic(err)
					}
				}
			}

			stockHotel, amount, err := p.agent.DetermineStockPurchase()
			if err != nil {
				panic(err)
			}

			if stockHotel > 0 && amount > 0 {
				err := p.buyStock(stockHotel, amount)
				if err != nil {
					panic(err)
				}
			}
		}

		newPlacedHotel := PlacedHotel{Pos: pos, Hotel: acquiringHotel}
		matrix.Set(pos.X, pos.Y, newPlacedHotel)

		// needs to propagate
		propagateHotelChain(p.game, newPlacedHotel)

		// only one effect can occur per placement
		return
	}

	// found a new chain - occurs when a tile has one or more neighbors which are all still undefined
	undefinedNeighbors := getUndefinedNeighbors(neighboringHotels)
	if len(undefinedNeighbors) > 0 {
		hotel, err := p.agent.DetermineHotelToFound()
		if err != nil {
			panic(err)
		}
		newPlacedHotel := PlacedHotel{Pos: pos, Hotel: hotel}
		matrix.Set(pos.X, pos.Y, newPlacedHotel)

		// miniature propagation of the hotel chain, these will never propagate beyond the orthogonal neighbors
		propagateHotelChain(p.game, newPlacedHotel)

		// receive a free stock in the chain for founding it
		_ = p.Inventory.Stocks[hotel].take(p.game.Inventory.Stocks[hotel])

		// only one effect can occur per placement
		return
	}

	panic("not sure how we got here, all situations should have been accounted for")
}

// anon func for the return syntax
func hasNeighboringHotel(neighbors []PlacedHotel) bool {
	for _, d := range util.Directions {
		if neighbors[d].Hotel != NoHotel {
			return true
		}
	}
	return false
}

func (p *Player) legalMoves() []Tile {
	return util.Filter(p.Inventory.Tiles.Items, func(val Tile) bool {
		return isLegalToPlace(p.game, val)
	})
}

func (p *Player) tiles() []Tile {
	return p.Inventory.Tiles.Items
}

// returnTile
// aliasing function for returning a specific tile to the game inventory
func (p *Player) returnTile(t Tile) {

	_, ok := p.Inventory.Tiles.indexOf(t)
	if !ok {
		return
	}

	_ = p.game.Inventory.Tiles.take(p.Inventory.Tiles)
}

// takeTile
// aliasing function for taking a tile from the game inventory
// this is normally a blind selection
func (p *Player) takeTile() {
	_ = p.Inventory.Tiles.take(p.game.Inventory.Tiles)
}

func (p *Player) takeTiles(other *Inventory, amount int) error {
	for i := 0; i < amount; i++ {
		err := p.Inventory.Tiles.take(other.Tiles)
		if err != nil {
			return err
		}
	}

	return nil
}

// refreshTiles
// when a player has no legal moves left to play, they can refresh their hand with this func
// puts all tiles back in the game inv, then takes 6 new ones
func (p *Player) refreshTiles() {

	// clone this so that it isn't getting fucked with as it removes items
	items := util.Clone(p.Inventory.Tiles.Items)
	for _, t := range items {
		p.returnTile(t)
	}

	// shuffle the tiles
	gameTiles := p.game.Inventory.Tiles.Items
	rand.Shuffle(len(gameTiles), func(i, j int) {
		gameTiles[i], gameTiles[j] = gameTiles[j], gameTiles[i]
	})

	for i := 0; i < 6; i++ {
		p.takeTile()
	}
}

// refreshOrSkip
// this function will refresh the tiles of a player if they have no legal moves repeatedly n times
// if the player doesn't have a valid move after a refresh then their turn should be skipped
// (as indicated by true in the returned bool)
func (p *Player) refreshOrSkip(n int) bool {
	for i := 0; i < n; i++ {
		if len(p.legalMoves()) < 1 {
			p.refreshTiles()
		} else {
			return false
		}
	}

	return true
}

func getUndefinedNeighbors(neighbors []PlacedHotel) []PlacedHotel {
	undefinedNeighbors := make([]PlacedHotel, 0)

	for _, d := range util.Directions {
		if neighbors[d].Hotel == UndefinedHotel {
			undefinedNeighbors = append(undefinedNeighbors, neighbors[d])
		}
	}

	return undefinedNeighbors
}

// getChainNeighbors
// returns a slice of the hotels around a position which are not NoHotel nor UndefinedHotel
func getChainNeighbors(neighbors []PlacedHotel) []PlacedHotel {
	chainNeighbors := make([]PlacedHotel, 0)

	for _, d := range util.Directions {
		if neighbors[d].Hotel != NoHotel && neighbors[d].Hotel != UndefinedHotel {
			chainNeighbors = append(chainNeighbors, neighbors[d])
		}
	}

	return chainNeighbors
}

// getChainsInNeighbors
// returns a slice of unique hotel chains which are in the neighbors slice
func getChainsInNeighbors(neighbors []PlacedHotel) []Hotel {
	chainNeighbors := getChainNeighbors(neighbors)
	hotels := util.Map(chainNeighbors, func(val PlacedHotel) Hotel {
		return val.Hotel
	})
	return util.UniqueElements[Hotel](hotels)
}

// getActiveHotelChains
// gets the active hotels, which are not on the board
// this is pretty brutally inefficient but whatever
func getActiveHotelChains(game *Game) []Hotel {
	chains := make([]Hotel, 0)

	game.Board.Matrix.Iterate(func(rt PlacedHotel, x int, y int, idx int) {
		chains = append(chains, rt.Hotel)
	})

	sortHotels(chains)

	return chains
}

// getActiveHotelChains
// gets the available hotels, which are not on the board
// this is pretty brutally inefficient but whatever
func getAvailableHotelChains(game *Game) []Hotel {
	chains := make(map[Hotel]struct{}, 0)

	for _, h := range HotelChainList {
		chains[h] = struct{}{}
	}

	game.Board.Matrix.Iterate(func(rt PlacedHotel, x int, y int, idx int) {
		delete(chains, rt.Hotel)
	})

	keys := util.Keys(chains)

	sortHotels(keys)

	return keys
}

func sortChainSizeMap(hotelSizeMap map[Hotel]int) []HotelChain {
	chainSizeSortedList := make([]HotelChain, 0)
	for hotel, size := range hotelSizeMap {
		chainSizeSortedList = append(chainSizeSortedList, HotelChain{
			Hotel: hotel,
			Size:  size,
		})
	}

	sort.Slice(chainSizeSortedList, func(i, j int) bool {
		return chainSizeSortedList[i].Size > chainSizeSortedList[j].Size
	})

	return chainSizeSortedList
}

func (p *Player) returnedAmount(tradeInAmount int) int {
	return tradeInAmount / 2
}

func (p *Player) canTradeIn(in Stock, out Stock, tradeInAmount int) error {

	if (tradeInAmount % 2) != 0 {
		return errors.New("trade in amount must be a multiple of two")
	}

	if p.remainingStock(in) < tradeInAmount {
		return errors.New("player does not have enough stock to trade in for")
	}

	returnAmount := p.returnedAmount(tradeInAmount)
	if returnAmount < 1 {
		return nil
	}

	outRemaining := p.game.remainingStock(out)

	if outRemaining < returnAmount {
		return errors.New("there is not enough stock to trade-in for")
	}

	return nil
}

func (p *Player) tradeIn(in Stock, out Stock, tradeInAmount int) error {

	err := p.canTradeIn(in, out, tradeInAmount)
	if err != nil {
		return err
	}

	returnAmount := p.returnedAmount(tradeInAmount)

	// return the stock to the bank
	_ = p.game.Inventory.takeHotelStock(Hotel(in), tradeInAmount, p.Inventory)

	// take half the amount of the new stock
	_ = p.Inventory.takeHotelStock(Hotel(out), returnAmount, p.game.Inventory)

	return nil
}

func (p *Player) canSellStock(stock Stock, amount int) error {
	if p.remainingStock(stock) < amount {
		return errors.New("amount sold cannot be greater than the amount the player has")
	}

	return nil
}

func (p *Player) sellStock(stock Stock, amount int) error {
	err := p.canSellStock(stock, amount)
	if err != nil {
		return err
	}

	hotel := Hotel(stock)
	err = p.game.Inventory.takeHotelStock(hotel, amount, p.Inventory)
	if err != nil {
		panic("shouldn't take more than there are")
	}

	chainSize := countHotelChain(p.game, hotel)
	value := sharesCalc(hotel, chainSize, amount)

	p.Inventory.Money += value

	return nil
}

func (p *Player) remainingStock(stock Stock) int {
	return len(p.Inventory.Stocks[Hotel(stock)].Items)
}
