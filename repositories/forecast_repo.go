package repositories

import "github.com/go-redis/redis"

type ForecastRepo struct {
	rc *redis.Conn
}

func NewForecastRepo(rc *redis.Conn) *ForecastRepo {
	return &ForecastRepo{
		rc: rc,
	}
}
