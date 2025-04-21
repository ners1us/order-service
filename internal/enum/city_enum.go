package enum

type City string

const (
	CityMoscow          City = "Москва"
	CitySaintPetersburg City = "Санкт-Петербург"
	CityKazan           City = "Казань"
)

func IsValidCity(city City) bool {
	switch city {
	case CityMoscow, CitySaintPetersburg, CityKazan:
		return true
	default:
		return false
	}
}

func (c City) String() string {
	return string(c)
}
