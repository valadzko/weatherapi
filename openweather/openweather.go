package openweather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/valadzko/weatherapi/models"
)

type OpenWeatherClient struct {
	Apikey  string
	BaseURL string
	cl      *http.Client
}

func NewOpenWeatherClient(a string) *OpenWeatherClient {
	owc := &OpenWeatherClient{
		Apikey:  a,
		BaseURL: "http://api.openweathermap.org/data/2.5",
		cl:      &http.Client{Timeout: 5 * time.Second},
	}
	return owc
}

func (owc *OpenWeatherClient) GetForecast(city, country string) (*models.Forecast, error) {
	url := fmt.Sprintf("%s/weather?q=%s,%s&appid=%s&units=metric", owc.BaseURL, city, country, owc.Apikey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("failed to reach open weather api")
	}
	defer resp.Body.Close()

	var f ApiForecast

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("failed to parse response body")
	}

	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Fatalln("failed unmarshal response body")
	}

	forecastModel := f.fromApiToModel()

	return forecastModel, nil
}

func (owc *OpenWeatherClient) GetForecastForDay(city, country string, day int) (*models.Forecast, error) {
	dayIndex := day + 1
	url := fmt.Sprintf("%s/forecast?q=%s,%s&cnt=%d&appid=%s&units=metric", owc.BaseURL, city, country, dayIndex, owc.Apikey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("failed to reach open weather api")
	}
	defer resp.Body.Close()

	var bf ApiBulkForecast

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("failed to parse response body")
	}

	err = json.Unmarshal(body, &bf)
	if err != nil {
		log.Fatalln("failed unmarshal response body")
	}
	if len(bf.ApiForecasts) < day {
		//return error
		log.Fatalln("there is not enough forecasts")
	}

	f := &bf.ApiForecasts[day]

	if bf.City != nil {
		f.Coord = &bf.City.Coord
		sys := Sys{
			Country: bf.City.Country,
			Sunrise: bf.City.Sunrise,
			Sunset:  bf.City.Sunset,
		}
		f.Sys = &sys
		f.Name = bf.City.Name
	}

	forecastModel := f.fromApiToModel()
	return forecastModel, nil
}

type ApiForecast struct {
	Coord      *Coord    `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       *Main     `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        *Sys      `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
	Pop        float64   `json:"pop"`
	Rain       Rain      `json:"rain"`
	DtTxt      string    `json:"dt_txt"`
}

func (api *ApiForecast) fromApiToModel() *models.Forecast {
	fm := models.Forecast{
		City:          api.Name,
		Wind:          "",
		Cloudiness:    "",
		RequestedTime: currentTimestampString(),
	}
	if api.Sys != nil {
		fm.Country = api.Sys.Country
		fm.Sunrise = daytime(api.Sys.Sunrise)
		fm.Sunset = daytime(api.Sys.Sunset)
	}
	if api.Main != nil {
		fm.Temperature = api.Main.Temp
		fm.Pressure = api.Main.Pressure
		fm.Humidity = api.Main.Humidity
	}
	if api.Coord != nil {
		fm.GeoCoordinates = fmt.Sprintf("[%.2f,%.2f]", api.Coord.Lon, api.Coord.Lat)
	}
	return &fm
}

func daytime(t int) string {
	timeT := time.Unix(int64(t), 0)
	return fmt.Sprintf("%d:%d", timeT.Hour(), timeT.Minute())
}

func currentTimestampString() string {
	t := time.Now()
	return t.Format("2006-02-03 09:15:05")
}

type ApiBulkForecast struct {
	Cod          string        `json:"cod"`
	Message      int           `json:"message"`
	Cnt          int           `json:"cnt"`
	ApiForecasts []ApiForecast `json:"list"`
	City         *City         `json:"city"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}
type Clouds struct {
	All int `json:"all"`
}
type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
	Pod     string `json:"pod"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
	TempKf    float64 `json:"temp_kf"`
}

type Rain struct {
	ThreeH float64 `json:"3h"`
}

type City struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Coord      Coord  `json:"coord"`
	Country    string `json:"country"`
	Population int    `json:"population"`
	Timezone   int    `json:"timezone"`
	Sunrise    int    `json:"sunrise"`
	Sunset     int    `json:"sunset"`
}
