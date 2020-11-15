package repositories

import "github.com/gomodule/redigo/redis"

type ForecastRepo struct {
	rc *redis.Conn
}

func NewForecastRepo(rc *redis.Conn) *ForecastRepo {
	return &ForecastRepo{
		rc: rc,
	}
}
