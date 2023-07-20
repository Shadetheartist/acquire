package acquire

type AIAgentSmart struct {
	player *Player
}

func aiAgentSmartFactory(player *Player) IAgent {
	return newAIAgentSmart(player)
}

func newAIAgentSmart(player *Player) *AIAgentSmart {
	return &AIAgentSmart{
		player: player,
	}
}

func (a *AIAgentSmart) DetermineTilePlacement() (Tile, error) {

	// does nothing if the player has legal moves to begin with
	skip := a.player.refreshOrSkip(1)
	if skip {
		return NoTile, nil
	}

	legalMoves := a.player.legalMoves()

	tile := legalMoves[0]

	return tile, nil
}

func (a *AIAgentSmart) DetermineHotelToFound() (Hotel, error) {
	chains := getAvailableHotelChains(a.player.game)
	return chains[0], nil
}

func (a *AIAgentSmart) DetermineHotelToMerge(hotels []Hotel) (Hotel, error) {
	return hotels[0], nil
}

func (a *AIAgentSmart) DetermineStockPurchase() (Hotel, int, error) {
	return NoHotel, 0, nil
}

func (a *AIAgentSmart) DetermineMergerAction(hotel Hotel) (MergerAction, error) {
	return Hold, nil
}

func (a *AIAgentSmart) DetermineTradeInAmount(acquiredHotel Hotel, acquiringHotel Hotel) (int, error) {
	return 0, nil
}

func (a *AIAgentSmart) DetermineStockSellAmount(hotel Hotel) (int, error) {
	return 0, nil
}

func (a *AIAgentSmart) DetermineGameEnd() (bool, error) {
	return true, nil
}
