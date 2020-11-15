package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/valadzko/weatherapi/models"
	"github.com/valadzko/weatherapi/openweather"
)

type WeatherHandler struct {
	repo models.ForecastRepository
	wc   *openweather.OpenWeatherClient
}

func NewWeatherHandler(repo models.ForecastRepository, wc *openweather.OpenWeatherClient) *WeatherHandler {
	return &WeatherHandler{
		repo: repo,
		wc:   wc,
	}
}

func (h *WeatherHandler) Weather(w http.ResponseWriter, r *http.Request) {
	request, err := parseRequest(*r.URL)
	if err != nil {
		// return
	}
	if err := request.validate(); err != nil {
		// return
	}

	// todo - searching for forecast in cache
	var f *models.Forecast

	// if not found - request via open weather api
	if request.day != nil {
		f, err = h.wc.GetForecastForDay(request.city, request.country, *request.day)
		if err != nil {
			// return
		}
	} else {
		f, err = h.wc.GetForecast(request.city, request.country)
		if err != nil {
			// return
		}
	}

	//todo save f to cache

	res := forecastModelToResponse(f)
	w.Header().Set("Content-type", "application/json")

	fmt.Fprint(w, res)
}

func forecastModelToResponse(f *models.Forecast) string {
	r := response{
		LocationName:   fmt.Sprintf("%s,%s", f.City, f.Country),
		Temperature:    fmt.Sprintf("%.2f °C", f.Temperature),
		Wind:           f.Wind,
		Cloudiness:     f.Cloudiness,
		Pressure:       fmt.Sprintf("%d hpa", f.Pressure),
		Humidity:       fmt.Sprintf("%d%%", f.Humidity),
		Sunrise:        f.Sunrise,
		Sunset:         f.Sunset,
		GeoCoordinates: f.GeoCoordinates,
		RequestedTime:  f.RequestedTime,
	}

	res, err := json.Marshal(r)
	if err != nil {
		log.Fatalln("failed to marshall response")
	}

	return string(res)
}

type response struct {
	LocationName   string `json:"location_name"`
	Temperature    string `json:"temperature"`
	Wind           string `json:"wind"`
	Cloudiness     string `json:"cloudiness"`
	Pressure       string `json:"pressure"`
	Humidity       string `json:"humidity"`
	Sunrise        string `json:"sunrise"`
	Sunset         string `json:"sunset"`
	GeoCoordinates string `json:"geo_coordinates"`
	RequestedTime  string `json:"requested_time"`
	Forecast       string `json:"forecast"`
}
type request struct {
	city    string
	country string
	day     *int
}

func (r *request) validate() error {
	if len(r.country) != 2 {
		return errors.New("country code must contain 2 characters")
	}

	if r.country != strings.ToLower(r.country) {
		// this could be just lowercased,
		// but I am leaving it like this according to the requirement
		// (Country is a country code of two characters in lowercase. Example: co)
		return errors.New("country parameter must contain only lowercase characters")
	}

	if len(r.city) <= 0 {
		return errors.New("city parameter can not be empty")
	}

	if r.day != nil {
		if *r.day < 0 || *r.day > 6 {
			return errors.New("day must be a number in [0,6]")
		}
	}
	return nil
}

func parseRequest(url url.URL) (request, error) {
	var res request
	params := url.Query()

	city, found := params["city"]
	if len(city) < 1 {
		// return error
		log.Fatalln("required city parameter is missed")
	}
	res.city = city[0]

	country, found := params["country"]
	if len(country) < 1 {
		// return error
		log.Fatalln("required country parameter is missed")
	}
	res.country = country[0]

	day, found := params["day"]
	if found {
		d, err := strconv.Atoi(day[0])
		if err != nil {
			//return error
			log.Fatalln("parameter day must be a number")
		}
		res.day = &d
	}

	return res, nil
}
