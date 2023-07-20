package acquire

type MergerAction int

const (
	Hold MergerAction = iota
	Trade
	Sell
)

type Stock Hotel

func (s Stock) String() string {
	return Hotel(s).String()
}

func (s Stock) getShareholders(game *Game) (*Player, int, *Player, int) {

	var majorShareholder *Player
	var majorShareholderShares int

	var minorShareholder *Player
	var minorShareholderShares int

	for _, p := range game.Players {
		h := Hotel(s)
		stocks := p.Inventory.Stocks[h].Items
		numShares := len(stocks)
		if numShares > majorShareholderShares {
			majorShareholder = p
			majorShareholderShares = numShares
		} else if numShares > minorShareholderShares {
			minorShareholder = p
			minorShareholderShares = numShares
		}
	}

	// if major shareholder is still null at this point, then there weren't any players holding shares in this chain
	if majorShareholder == nil {
		return nil, 0, nil, 0
	}

	// if there is no minor shareholder, the major shareholder becomes the major AND minor shareholder
	if minorShareholder == nil {
		minorShareholder = majorShareholder
		minorShareholderShares = majorShareholderShares
	}

	return majorShareholder, majorShareholderShares, minorShareholder, minorShareholderShares
}
