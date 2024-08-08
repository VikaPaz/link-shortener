package app

import (
	"database/sql"
	"github.com/joho/godotenv"
	rConn "github.com/redis/go-redis/v9"
	"links-shorter/internal/server"
	"links-shorter/internal/service"
	"links-shorter/repository/postgres"
	"links-shorter/repository/redis"
	"log"
	"net/http"
	"os"
	"time"
)

const upt = 10 * time.Millisecond

type PostgresConfig struct {
	DB postgres.Config
}

type RedisConfig struct {
	DB redis.Config
}

func Run() error {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	var (
		dbConn    *sql.DB
		cacheConn *rConn.Client
		err       error
	)

	confPostgers := PostgresConfig{
		DB: postgres.Config{
			Host:     os.Getenv("HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("USER"),
			Password: os.Getenv("PASSWORD"),
			Dbname:   os.Getenv("DB_NAME"),
		},
	}

	confRedis := RedisConfig{
		DB: redis.Config{
			Host:     os.Getenv("HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("PASSWORD"),
		},
	}

	dbConn, err = postgres.Connection(confPostgers.DB)
	if err != nil {
		return err
	}
	repo := postgres.NewRepository(dbConn)

	cacheConn, err = redis.Connection(confRedis.DB)
	if err != nil {
		return err
	}
	cache := redis.NewRepository(cacheConn)

	service := service.NewService(repo, cache)

	out := make(chan int)
	go service.Writer(upt, out)

	err = http.ListenAndServe(":3000", server.NewServer(service).Handlers())
	return err
}
