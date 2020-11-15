package repositories

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/valadzko/weatherapi/models"
)

type ForecastRepo struct {
	rdb redis.Client
}

func NewForecastRepo(rdb redis.Client) *ForecastRepo {
	return &ForecastRepo{
		rdb: rdb,
	}
}

func (fr *ForecastRepo) Save(f *models.Forecast) error {

	key := fmt.Sprintf("%s:%s", f.City, f.Country)
	_ = key

	err := rdb.Set(ctx, "key1", "value1", 0).Err()
	if err != nil {
		panic(err)
	}

	return nil
}
