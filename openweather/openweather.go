package openweather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

func (owc *OpenWeatherClient) GetForecast(city, country string) *Forecast {
	url := fmt.Sprintf("%s/weather?q=%s,%s&appid=%s&units=metric", owc.BaseURL, city, country, owc.Apikey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("failed to reach open weather api")
	}
	defer resp.Body.Close()

	var f Forecast

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("failed to parse response body")
	}

	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Fatalln("failed unmarshal response body")
	}

	return &f
}

func (owc *OpenWeatherClient) GetForecastForDay(city, country string, day int) *Forecast {
	dayIndex := day - 1
	url := fmt.Sprintf("%s/forecast?q=%s,%s&cnt=%d&appid=%s&units=metric", owc.BaseURL, city, country, dayIndex, owc.Apikey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("failed to reach open weather api")
	}
	defer resp.Body.Close()

	var bf BulkForecast

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("failed to parse response body")
	}

	err = json.Unmarshal(body, &bf)
	if err != nil {
		log.Fatalln("failed unmarshal response body")
	}

	if len(bf.Forecasts) < day {
		log.Fatalln("there is not enough forecasts")
	}

	f := &bf.Forecasts[dayIndex]

	f.Coord = bf.City.Coord

	sys := Sys{
		Country: bf.City.Country,
		Sunrise: bf.City.Sunrise,
		Sunset:  bf.City.Sunset,
	}
	f.Sys = sys
	f.Name = bf.City.Name

	return f
}

type Forecast struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
	Pop        float64   `json:"pop"`
	Rain       Rain      `json:"rain"`
	DtTxt      string    `json:"dt_txt"`
}

type BulkForecast struct {
	Cod       string     `json:"cod"`
	Message   int        `json:"message"`
	Cnt       int        `json:"cnt"`
	Forecasts []Forecast `json:"list"`
	City      City       `json:"city"`
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
