package models

type Forecast struct {
	City           string
	Country        string
	Temperature    float64
	Wind           string
	Cloudiness     string
	Pressure       int
	Humidity       int
	Sunrise        string
	Sunset         string
	GeoCoordinates string
	RequestedTime  string
}

type ForecastRepository interface {
	Save(f *Forecast) error
	SaveWithDay(f *Forecast, day int) error
	FindByCityAndCountry(city, country string) (*Forecast, error)
	FindByCityCountryAndDay(city, country string, day int) (*Forecast, error))
}
