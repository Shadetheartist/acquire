package acquire

// AIAgentStupid
// this Agent makes moves without thinking at all.
// basically playing the first legal move it sees
type AIAgentStupid struct {
	player *Player
}

func aiAgentStupidFactory(player *Player) IAgent {
	return newAIAgentStupid(player)
}

func newAIAgentStupid(player *Player) *AIAgentStupid {
	return &AIAgentStupid{
		player: player,
	}
}

func (a *AIAgentStupid) DetermineTilePlacement() (Tile, error) {

	// does nothing if the player has legal moves to begin with
	skip := a.player.refreshOrSkip(1)
	if skip {
		return NoTile, nil
	}

	legalMoves := a.player.legalMoves()

	tile := legalMoves[0]

	return tile, nil
}

func (a *AIAgentStupid) DetermineHotelToFound() (Hotel, error) {
	chains := getAvailableHotelChains(a.player.game)
	return chains[0], nil
}

func (a *AIAgentStupid) DetermineHotelToMerge(hotels []Hotel) (Hotel, error) {
	return hotels[0], nil
}

func (a *AIAgentStupid) DetermineStockPurchase() (Hotel, int, error) {
	return NoHotel, 0, nil
}

func (a *AIAgentStupid) DetermineMergerAction(hotel Hotel) (MergerAction, error) {
	return Hold, nil
}

func (a *AIAgentStupid) DetermineTradeInAmount(acquiredHotel Hotel, acquiringHotel Hotel) (int, error) {
	return 0, nil
}

func (a *AIAgentStupid) DetermineStockSellAmount(hotel Hotel) (int, error) {
	return 0, nil
}

func (a *AIAgentStupid) DetermineGameEnd() (bool, error) {
	return true, nil
}
