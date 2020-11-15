package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/valadzko/weatherapi/models"
)

type ForecastRepo struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewForecastRepo(rdb *redis.Client, ttl time.Duration) *ForecastRepo {
	return &ForecastRepo{
		rdb: rdb,
		ttl: ttl,
	}
}

func (fr *ForecastRepo) Save(f *models.Forecast) error {
	key := fmt.Sprintf("%s:%s", f.City, f.Country)
	return fr.save(key, f)
}

func (fr *ForecastRepo) SaveWithDay(f *models.Forecast, day int) error {
	key := fmt.Sprintf("%s:%s:%d", f.City, f.Country, day)
	return fr.save(key, f)
}

func (fr *ForecastRepo) FindByCityAndCountry(city, country string) (*models.Forecast, error) {
	key := fmt.Sprintf("%s:%s", city, country)
	return fr.find(key)
}

func (fr *ForecastRepo) FindByCityCountryAndDay(city, country string, day int) (*models.Forecast, error) {
	key := fmt.Sprintf("%s:%s:%d", city, country, day)
	return fr.find(key)
}

func (fr *ForecastRepo) save(key string, f *models.Forecast) error {
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

func (fr *ForecastRepo) find(key string) (*models.Forecast, error) {
	ctx := context.Background()
	val, err := fr.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		panic(err)
	}

	var f *models.Forecast

	err = json.Unmarshal([]byte(val), &f)
	if err != nil {
		fmt.Println("Could not unmarshal forecast from cache")
	}

	return f, nil
}
