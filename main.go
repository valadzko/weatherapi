package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/valadzko/weatherapi/openweather"
)

var (
	APIKEY = "1508a9a4840a5574c822d70ca2132032"
	PORT   = ":8080"
)

func main() {
	http.HandleFunc("/weather", weatherHandler)

	fmt.Println("Started server at port 8080:")
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
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

	f := owc.GetForecast(city[0], country[0])
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
	temperature := fmt.Sprintf("%.2f °C", kelvinToCelcius(f.Main.Temp))
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

func kelvinToCelcius(k float64) float64 {
	return k - 273.15
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
