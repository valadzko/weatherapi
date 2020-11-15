package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/davecgh/go-spew/spew"
	"github.com/valadzko/weatherapi/controllers"
	"github.com/valadzko/weatherapi/openweather"
)

func main() {
	//read APIKEY from ENV
	apikey := ""

	//read application port from ENV
	port := "8080"

	// connect to redis
	rc, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	// create repo
	repo := NewForecastRepo(rc)

	// create open weather client
	owc := openweather.NewOpenWeatherClient(apikey)

	// create handler
	h := controllers.NewWeatherHandler(repo, owc)
	http.HandleFunc("/weather", h.Weather)

	// create and start server
	s := &http.Server{Addr: fmt.Sprintf("127.0.0.1:%s", port)}
	s.ListenAndServe()
}

func weatherHandler(w http.ResponseWriter, req *http.Request) {

	owc := openweather.NewOpenWeatherClient(APIKEY)

	params := req.URL.Query()
	city, found := params["city"]
	if !found || len(city) < 1 {
		log.Fatalln("city param is missed")
	}
	country, found := params["country"]
	if !found || len(country) < 1 {
		log.Fatalln("country param is missed")
	}

	var f *openweather.Forecast

	day, found := params["day"]
	if found {
		d, _ := strconv.Atoi(day[0])
		f = owc.GetForecastForDay(city[0], country[0], d)
	} else {
		f = owc.GetForecast(city[0], country[0])
	}

	resp := forecastToResponse(f)

	r, err := json.Marshal(resp)
	if err != nil {
		log.Fatalln("failed to marshall response")
	}

	spew.Dump(resp)
	w.Header().Set("Content-type", "application/json")

	fmt.Fprint(w, string(r))
}

func forecastToResponse(f *openweather.Forecast) *response {
	r := &response{}
	name := fmt.Sprintf("%s,%s", f.Name, f.Sys.Country)
	temperature := fmt.Sprintf("%.2f °C", f.Main.Temp)
	// wind =
	// cloudiness :=
	pressure := fmt.Sprintf("%d hpa", f.Main.Pressure)
	humidity := fmt.Sprintf("%d%%", f.Main.Humidity)
	sunrise := daytime(f.Sys.Sunrise)
	sunset := daytime(f.Sys.Sunset)
	geo := fmt.Sprintf("[%.2f,%.2f]", f.Coord.Lon, f.Coord.Lat)

	r.LocationName = name
	r.Temperature = temperature
	// r.Wind = wind
	// r.Cloudiness = cloudiness
	r.Pressure = pressure
	r.Humidity = humidity
	r.Sunrise = sunrise
	r.Sunset = sunset
	r.GeoCoordinates = geo
	r.RequestedTime = currentTimestampString()
	return r
}

func daytime(t int) string {
	timeT := time.Unix(int64(t), 0)
	return fmt.Sprintf("%d:%d", timeT.Hour(), timeT.Minute())
}

func currentTimestampString() string {
	t := time.Now()
	return t.Format("2006-02-03 09:15:05")
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
