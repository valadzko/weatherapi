package models

type Forecast struct {
	City           string
	Country        string
	Temperature    string
	Wind           string
	Cloudiness     string
	Pressure       string
	Humidity       string
	Sunrise        string
	GeoCoordinates string
	RequestedTime  string
}

type ForecastRepository interface {
	FindByCityAndCountry(city, country string) (*Forecast, error)
	Save(f *Forecast) error
}
