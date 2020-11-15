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
	var err error
	request, err := parseRequest(*r.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - Bad request:%s", err.Error())))
		return
	}
	if err := request.validate(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(fmt.Sprintf("422 - Unprocessable Entity:%s", err.Error())))
		return
	}
	var f *models.Forecast

	if request.day == nil {
		fmt.Printf("GET /weather?&city=%s,country=%s\n", request.city, request.country)

		// lookup in cache
		f, _ = h.repo.FindByCityAndCountry(request.city, request.country)
		// request forecast from server
		if f == nil {
			f, err = h.wc.GetForecast(request.city, request.country)
			if err != nil {
				// could not find forecast in both cache and remote server
				fmt.Printf("could not find forecast: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("500 - Something bad happened!")))
				return
			}
			if f != nil {
				// caching forecast
				if err := h.repo.Save(f); err != nil {
					fmt.Printf("failed to cache forecast: %d\n", err)
				}
			}
		}
	} else {
		// day parameter is present
		fmt.Printf("GET /weather?&city=%s&country=%s&day=%s\n", request.city, request.country, *request.day)
		// lookup in cache
		f, _ = h.repo.FindByCityCountryAndDay(request.city, request.country, *request.day)
		// request forecast from server
		if f == nil {
			f, err = h.wc.GetForecastForDay(request.city, request.country, *request.day)
			if err != nil {
				// could not find forecast in both cache and remote server
				fmt.Printf("could not find forecast: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("500 - Something bad happened!")))
				return
			}
			if f != nil {
				// caching forecast
				if err := h.repo.SaveWithDay(f, *request.day); err != nil {
					fmt.Printf("failed to cache forecast: %d\n", err)
				}
			}
		}
	}

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
		// but I am leaving it like this according to the requirement:
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
		return res, errors.New("required city parameter is missed")
	}
	res.city = city[0]

	country, found := params["country"]
	if len(country) < 1 {
		return res, errors.New("required country parameter is missed")
	}
	res.country = country[0]

	day, found := params["day"]
	if found {
		d, err := strconv.Atoi(day[0])
		if err != nil {
			return res, errors.New("parameter day must be a number")
		}
		res.day = &d
	}

	return res, nil
}
