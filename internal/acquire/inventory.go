package acquire

type INamed interface {
	Name() string
}

type Inventory struct {
	Owner  INamed
	Money  int
	Tiles  *Collection[Tile]
	Stocks map[Hotel]*Collection[Stock]
}

func newInventory(owner INamed, money int) *Inventory {
	stocks := make(map[Hotel]*Collection[Stock])

	for _, h := range HotelChainList {
		stocks[h] = newCollection[Stock]()
	}

	return &Inventory{
		Owner:  owner,
		Tiles:  newCollection[Tile](),
		Money:  money,
		Stocks: stocks,
	}
}

// takeMoney
// takes money from another inventory, if they don't have enough, it takes as much as it can
// this function returns the amount of money actually taken
func (inv *Inventory) takeMoney(other *Inventory, amount int) int {
	// if we would take more than they have, just take everything and leve them with zero
	if amount > other.Money {
		amount = other.Money
	}

	other.Money -= amount
	inv.Money += amount

	return amount
}

func (inv *Inventory) takeHotelStock(hotel Hotel, amount int, other *Inventory) error {
	for i := 0; i < amount; i++ {
		err := inv.Stocks[hotel].take(other.Stocks[hotel])
		if err != nil {
			return err
		}
	}

	return nil
}
