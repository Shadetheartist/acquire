package acquire_2

import "sort"

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

func hotelsAsInitials(hotels []Hotel) []string {
	initials := make([]string, len(hotels))
	for i, h := range hotels {
		initials[i] = hotelInitials[h]
	}
	return initials
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

var HotelList = []Hotel{
	NoHotel,
	UndefinedHotel,
	WorldwideHotel,
	SacksonHotel,
	FestivalHotel,
	ImperialHotel,
	AmericanHotel,
	ContinentalHotel,
	TowerHotel,
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

func sortHotels(hotels []Hotel) {
	sort.Slice(hotels, func(i, j int) bool {
		return hotels[i] < hotels[j]
	})
}
