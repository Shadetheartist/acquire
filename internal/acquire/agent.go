package acquire

import (
	"acquire/internal/util"
	"errors"
)

type IAgent interface {
	DetermineTilePlacement() (Tile, error)
	DetermineHotelToFound() (Hotel, error)
	DetermineHotelToMerge(hotels []Hotel) (Hotel, error)
	DetermineStockPurchase() (Hotel, int, error)
	DetermineMergerAction(hotel Hotel) (MergerAction, error)
	DetermineTradeInAmount(acquiredHotel Hotel, acquiringHotel Hotel) (int, error)
	DetermineStockSellAmount(hotel Hotel) (int, error)
	DetermineGameEnd() (bool, error)
}

type HumanAgent struct {
	player *Player
}

func humanAgentFactory(player *Player) IAgent {
	return newHumanAgent(player)
}

func newHumanAgent(player *Player) *HumanAgent {
	return &HumanAgent{
		player: player,
	}
}

func (a *HumanAgent) DetermineHotelToFound() (Hotel, error) {
	hotel := getHotelInput(a.player.game, "Founding new Chain", getAvailableHotelChains(a.player.game))
	return hotel, nil
}

func (a *HumanAgent) DetermineHotelToMerge(hotels []Hotel) (Hotel, error) {
	hotel := getHotelInput(a.player.game, "Merger", hotels)
	return hotel, nil
}

func (a *HumanAgent) DetermineTilePlacement() (Tile, error) {
	if len(a.player.legalMoves()) < 1 {
		return NoTile, errors.New("no legal moves left to play")
	}

	tile := getTileInput(a.player.game, a.player)

	return tile, nil
}

func (a *HumanAgent) DetermineStockPurchase() (Hotel, int, error) {
	availableStocks := a.player.game.remainingStocks()
	hotel := getBuyStockInput(a.player.game)
	if hotel == NoHotel {
		return NoHotel, 0, nil
	}

	stocksAvailable := availableStocks[Stock(hotel)]
	costPerShare := sharesCalc(hotel, countHotelChain(a.player.game, hotel), 1)

	canAfford := a.player.Inventory.Money / costPerShare
	validAmount := util.Min(stocksAvailable, canAfford)

	amount := getNumStocksToBuy(a.player.game, hotel, util.Min(validAmount, 3))

	return hotel, amount, nil
}

func (a *HumanAgent) DetermineMergerAction(hotel Hotel) (MergerAction, error) {
	return getMergerActionTypeInput(a.player.game, hotel), nil
}

func (a *HumanAgent) DetermineTradeInAmount(acquiredHotel Hotel, acquiringHotel Hotel) (int, error) {
	return getTradeInAmount(a.player.game, a.player, Stock(acquiredHotel), Stock(acquiringHotel)), nil
}

func (a *HumanAgent) DetermineStockSellAmount(hotel Hotel) (int, error) {
	return getSellAmount(a.player.game, a.player, Stock(hotel)), nil
}

func (a *HumanAgent) DetermineGameEnd() (bool, error) {
	// - After playing a tile, the player whose turn it is declares that there is a hotel chain of size 41 or more on the game board.
	// - After playing a tile, the player whose turn it is declares that all the hotel chains on the game board are safe.
	return true, nil
}

type AIAgent struct {
	player *Player
}

func aiAgentFactory(player *Player) IAgent {
	return newAIAgent(player)
}

func newAIAgent(player *Player) *AIAgent {
	return &AIAgent{
		player: player,
	}
}

func (a *AIAgent) DetermineTilePlacement() (Tile, error) {

	// does nothing if the player has legal moves to begin with
	skip := a.player.refreshOrSkip(1)
	if skip {
		return NoTile, nil
	}

	legalMoves := a.player.legalMoves()

	tile := legalMoves[0]

	return tile, nil
}

func (a *AIAgent) DetermineHotelToFound() (Hotel, error) {
	chains := getAvailableHotelChains(a.player.game)
	return chains[0], nil
}

func (a *AIAgent) DetermineHotelToMerge(hotels []Hotel) (Hotel, error) {
	return hotels[0], nil
}

func (a *AIAgent) DetermineStockPurchase() (Hotel, int, error) {
	return NoHotel, 0, nil
}

func (a *AIAgent) DetermineMergerAction(hotel Hotel) (MergerAction, error) {
	return Hold, nil
}

func (a *AIAgent) DetermineTradeInAmount(acquiredHotel Hotel, acquiringHotel Hotel) (int, error) {
	return 0, nil
}

func (a *AIAgent) DetermineStockSellAmount(hotel Hotel) (int, error) {
	return 0, nil
}

func (a *AIAgent) DetermineGameEnd() (bool, error) {
	return true, nil
}
