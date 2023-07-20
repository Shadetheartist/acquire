package acquire

import "sort"

type Hotel int

type HotelChain struct {
	Hotel Hotel
	Size  int
}

func (h Hotel) String() string {
	return hotelNames[h]
}

func (h Hotel) Initial() string {
	return hotelInitials[h]
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

func isActualHotelChain(hotel Hotel) bool {
	for _, h := range HotelChainList {
		if h == hotel {
			return true
		}
	}

	return false
}

func hotelFromInitial(str string) Hotel {
	for idx, initial := range hotelInitials {
		if initial == str {
			return HotelList[idx]
		}
	}

	return NoHotel
}

func hotelsWithPurchasableStock(game *Game) []Hotel {
	hotels := make([]Hotel, 0)

	stocks := game.purchasableStocks()
	for s, n := range stocks {
		if n > 0 {
			hotels = append(hotels, Hotel(s))
		}
	}

	sortHotels(hotels)

	return hotels
}

func sortHotels(hotels []Hotel) {
	sort.Slice(hotels, func(i, j int) bool {
		return hotels[i] < hotels[j]
	})
}
