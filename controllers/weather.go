package controllers

import (
	"errors"
	"http"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/valadzko/weatherapi/models"
	"github.com/valadzko/weatherapi/openweather"
)

type request struct {
	city    string
	country string
	day     *int
}

func (r *request) validate() error {
	if len(r.country) != 2 {
		return errors.New("country code must contain 2 characters")
	}

	if r.country != strings.ToLower() {
		// this could be just lowercased, but leaving it like this according to the requirement
		return errors.New("country code must contain only lowercase characters")
	}

	if len(city) <= 0 {
		return errors.New("city parameter can not be empty")
	}

	if day != nil && (day < 0 || day > 6) {
		return errors.New("day must be a number in [0,6]")
	}
	return nil
}

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

func (h *WeatherHandler) Weather(w http.ResponseWriter, r *http.request) {
	request, err := parseRequest(r.URL)
	if err != nil {
		// return
	}
	if err := request.validate(); err != nil {
		// return
	}

	//todo - search in repo

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
