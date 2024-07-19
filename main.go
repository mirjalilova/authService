package main

import (
	"context"
	"log"
	"net"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	"golang.org/x/exp/slog"

	"github.com/mirjalilova/authService/api"
	_ "github.com/mirjalilova/authService/api/docs"
	"github.com/mirjalilova/authService/api/handlers"
	"github.com/mirjalilova/authService/config"
	kafka "github.com/mirjalilova/authService/consumer"
	pb "github.com/mirjalilova/authService/genproto/auth"
	"github.com/mirjalilova/authService/service"
	"github.com/mirjalilova/authService/storage/postgres"
)

func main() {
	cfg := config.Load()

	db, err := postgres.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatalf("can't connect to db: %v", err)
	}
	defer db.Db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	brokers := []string{"kafka:9092"}

	kcm := kafka.NewKafkaConsumerManager()
	authService := service.NewAuthService(db)
	userService := service.NewUserService(db)

	if err := kcm.RegisterConsumer(brokers, "reg-user", "auth", kafka.UserRegisterHandler(authService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			slog.Warn("Consumer for topic 'reg-user' already exists")
		} else {
			slog.Error("Error registering consumer: %v", err)
		}
	}
	if err := kcm.RegisterConsumer(brokers, "upd-user", "auth", kafka.UserEditProfileHandler(userService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			slog.Warn("Consumer for topic 'upd-user' already exists")
		} else {
			slog.Error("Error registering consumer: %v", err)
		}
	}
	if err := kcm.RegisterConsumer(brokers, "upd-pass", "auth", kafka.UserEditPasswordHandler(userService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			slog.Warn("Consumer for topic 'upd-pass' already exists")
		} else {
			slog.Error("Error registering consumer: %v", err)
		}
	}
	if err := kcm.RegisterConsumer(brokers, "upd-setting", "auth", kafka.UserEditSettingHandler(userService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			slog.Warn("Consumer for topic 'upd-setting' already exists")
		} else {
			slog.Error("Error registering consumer: %v", err)
		}
	}

	listener, err := net.Listen("tcp", cfg.USER_PORT)
	if err != nil {
		slog.Error("can't listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)

	go func() {
		slog.Info("gRPC server started on port %s", cfg.USER_PORT)
		if err := s.Serve(listener); err != nil {
			slog.Error("can't serve: %v", err)
		}
	}()

	h := handlers.NewHandler(authService, rdb)

	router := api.Engine(h)
	if err := router.Run(cfg.AUTH_PORT); err != nil {
		log.Fatalf("can't start server: %v", err)
	}

	log.Printf("REST server started on port %s", cfg.AUTH_PORT)
}
