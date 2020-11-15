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
	//	FindByCityAndCountry(city, country string) (*Forecast, error)
	Save(f *Forecast) error
}
