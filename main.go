package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"

	"github.com/mirjalilova/authService/api"
	_ "github.com/mirjalilova/authService/api/docs"
	"github.com/mirjalilova/authService/api/handlers"
	"github.com/mirjalilova/authService/config"
	"github.com/mirjalilova/authService/service"
	"github.com/mirjalilova/authService/storage/postgres"
	kafka "github.com/mirjalilova/authService/consumer"
)

func main() {
	cfg := config.Load()

	db, err := postgres.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatalf("can't connect to db: %v", err)
	}
	defer db.Db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	brokers := []string{"localhost:9092"}

	kcm := kafka.NewKafkaConsumerManager()
	authService := service.NewAuthService(db)

	if err := kcm.RegisterConsumer(brokers, "reg-user", "auth", kafka.UserRegisterHandler(authService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'reg-user' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}

	h := handlers.NewHandler(authService, rdb)

	router := api.Engine(h)
	if err := router.Run(cfg.AUTH_PORT); err != nil {
		log.Fatalf("can't start server: %v", err)
	}

	log.Printf("Server started on port %s", cfg.AUTH_PORT)
}
