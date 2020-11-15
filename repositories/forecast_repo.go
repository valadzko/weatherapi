package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/valadzko/weatherapi/models"
)

type ForecastRepo struct {
	rdb redis.Client
	ttl int
}

func NewForecastRepo(rdb redis.Client, ttl int) *ForecastRepo {
	return &ForecastRepo{
		rdb: rdb,
		ttl: ttl,
	}
}

func (fr *ForecastRepo) Save(f *models.Forecast) error {
	key := fmt.Sprintf("%s:%s", f.City, f.Country)
	return fr.save(key, f)
}

func (fr *ForecastRepo) SaveWithDay(f *models.Forecast, int day) error {
	key := fmt.Sprintf("%s:%s:%d", f.City, f.Country, day)
	return fr.save(key, f)
}

func (fr *ForecastRepo) FindByCityAndCountry(city, country string) (*Forecast, error) {
	key := fmt.Sprintf("%s:%s", f.City, f.Country)
	return fr.find(key)
}

func (fr *ForecastRepo) FindByCityCountryAndDay(city, country string, day int) (*Forecast, error) {
	key := fmt.Sprintf("%s:%s:%d", f.City, f.Country, day)
	return fr.find(key)
}

func (fr *ForecastRepo) save(key string, f *Forecast) error {
	ctx := context.Background()
	bytes, err := json.Marshal(f)
	if err != nil {
		fmt.Println("failed to marshal forecast")
	}

	err = fr.rdb.Set(ctx, key, string(bytes), fr.ttl).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (fr *ForecastRepo) find(key string) (*Forecast, error) {
	val, err := fr.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	}

	var f *Forecast

	err := json.Unmarshal(val, &f)
	if err != nil {
		fmt.Println("Could not unmarshal forecast from cache")
	}

	return f, nil
}
