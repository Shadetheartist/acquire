package acquire

import "errors"

type Hotel int

func (h Hotel) String() string {
	return hotelNames[h]
}

func (h Hotel) Initial() string {
	return hotelInitials[h]
}

// Index
// index of hotel chain in the HotelChainList
func (h Hotel) Index() int {
	// -1 for no hotel, -1 for undefined hotel = -2
	idx := int(h) - 2

	if idx < 0 {
		panic("hotel index used for NoHotel or UndefinedHotel (or other?)")
	}

	return idx
}

var hotelInitials = []string{
	"",  //NoHotel
	"=", //UndefinedHotel
	"W", //WorldwideHotel
	"S", //SacksonHotel
	"F", //FestivalHotel
	"I", //ImperialHotel
	"A", //AmericanHotel
	"C", //ContinentalHotel
	"T", //TowerHotel
}

var hotelNames = []string{
	"No Hotel",
	"Undefined",
	"Worldwide",
	"Sackson",
	"Festival",
	"Imperial",
	"American",
	"Continental",
	"Tower",
}

// HotelChainList
// similar to HotelList, but only actual valid chains
var HotelChainList = []Hotel{
	WorldwideHotel,
	SacksonHotel,
	FestivalHotel,
	ImperialHotel,
	AmericanHotel,
	ContinentalHotel,
	TowerHotel,
}

const (
	NoHotel Hotel = iota
	UndefinedHotel
	WorldwideHotel
	SacksonHotel
	FestivalHotel
	ImperialHotel
	AmericanHotel
	ContinentalHotel
	TowerHotel
)

// ChainFromIdx
// maps the index range [0-6] to valid hotel chains from HotelChainList
func ChainFromIdx(idx int) Hotel {
	return HotelChainList[idx]
}

// ChainFromInitial
// Returns the hotel chain associated with an initial, or otherwise NoHotel
func ChainFromInitial(initial string) (Hotel, error) {
	for idx, s := range hotelInitials {
		if s == initial {
			return Hotel(idx), nil
		}
	}

	return NoHotel, errors.New("chain not found from initial " + initial)
}

func (h Hotel) Value(game *Game, amount int) int {
	size := game.ChainSize[h.Index()]
	return sharesCalc(h, size, amount)
}
