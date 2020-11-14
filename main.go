package main

import (
	"fmt"
	"log"
	"net/http"

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

	spew.Dump(f)

	fmt.Fprintf(w, "test response for /weather")
}
