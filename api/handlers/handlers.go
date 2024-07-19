package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/mirjalilova/authService/producer"
	"github.com/mirjalilova/authService/service"
	"log"
)

type Handlers struct {
	Auth     *service.AuthService
	RDB      *redis.Client
	Producer kafka.KafkaProducer
}

func NewHandler(auth *service.AuthService, rdb *redis.Client) *Handlers {
	pr, err := kafka.NewKafkaProducer([]string{"kafka:9092"})
	if err != nil {
		log.Fatal(err)
	}
	return &Handlers{
		Auth:     auth,
		RDB:      rdb,
		Producer: pr,
	}
}
