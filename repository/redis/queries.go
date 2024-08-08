package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type repositoryImpl struct {
	conn *redis.Client
}

func (r *repositoryImpl) Get(link string) (string, error) {
	ctx := context.Background()

	var (
		v   string
		err error
	)
	_, err = r.conn.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		v, err = pipeliner.Get(ctx, link).Result()
		if err != nil {
			return err
		}
		if v != "" {
			pipeliner.Expire(context.Background(), link, 10*time.Minute)
		}

		return nil
	})
	if err == redis.Nil {
		return v, nil
	}

	return v, err
}

func (r *repositoryImpl) Set(m map[string]string) error {
	ctx := context.Background()

	_, err := r.conn.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for key, value := range m {
			pipeliner.Set(ctx, key, value, 10*time.Minute)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
