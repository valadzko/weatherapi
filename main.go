package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/valadzko/weatherapi/controllers"
	"github.com/valadzko/weatherapi/openweather"
	"github.com/valadzko/weatherapi/repositories"
)

func main() {

	apikey := getEnv("APIKEY", "1508a9a4840a5574c822d70ca2132032")
	port := getEnv("PORT", "8080")
	redisHost := getEnv("REDIS_HOST", "127.0.0.1")
	redisPort := getEnv("REDIS_PORT", "6379")
	ttl := 2 * time.Minute

	// init and connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password:   "", // no password set
		DB:         0,  // use default DB
		MaxRetries: 3,
	})
	ctx := context.Background()
	pong, _ := rdb.Ping(ctx).Result()
	if pong != "PONG" {
		log.Fatalf("could not connect to redis: %s:%s", redisHost, redisPort)
	}
	defer rdb.Close()

	// create repo
	repo := repositories.NewForecastRepo(rdb, ttl)

	// create open weather client
	owc := openweather.NewOpenWeatherClient(apikey)

	// create handler
	h := controllers.NewWeatherHandler(repo, owc)
	http.HandleFunc("/weather", h.Weather)

	// create and start server
	s := &http.Server{Addr: fmt.Sprintf("127.0.0.1:%s", port)}
	fmt.Printf("Started server at port :%s\n", port)
	log.Fatal(s.ListenAndServe())
}

func getEnv(key, defvalue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defvalue
	}
	return value
}
