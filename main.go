package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"

	"github.com/valadzko/weatherapi/controllers"
	"github.com/valadzko/weatherapi/openweather"
	"github.com/valadzko/weatherapi/repositories"
)

func main() {

	apikey := getEnv("APIKEY", "1508a9a4840a5574c822d70ca2132032")
	port := getEnv("PORT", "8080")
	redisHost := getEnv("REDIS_HOST", "127.0.0.1")
	redisPort := getEnv("REDIS_PORT", "6379")

	// connect to redis
	rc, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort))
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	// create repo
	repo := repositories.NewForecastRepo(nil)

	// create open weather client
	owc := openweather.NewOpenWeatherClient(apikey)

	// create handler
	h := controllers.NewWeatherHandler(repo, owc)
	http.HandleFunc("/weather", h.Weather)

	// create and start server
	s := &http.Server{Addr: fmt.Sprintf("127.0.0.1:%s", port)}
	fmt.Printf("Started server at port :%s\n", port)
	s.ListenAndServe()
}

func getEnv(key, defvalue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defvalue
	}
	return value
}
