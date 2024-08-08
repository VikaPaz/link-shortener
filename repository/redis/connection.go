package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Config struct{ Host, Port, Password string }

func Connection(conf Config) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: conf.Password,
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewRepository(conn *redis.Client) *repositoryImpl {
	return &repositoryImpl{
		conn: conn,
	}
}
