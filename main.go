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
	f := owc.GetForecast("Bogota", "co")

	spew.Dump(f)

	fmt.Fprintf(w, "test response for /weather")
}
