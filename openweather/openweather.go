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
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
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
}
