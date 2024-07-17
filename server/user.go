package main

import (
	"log"
	"net"

	"golang.org/x/exp/slog"

	cf "github.com/mirjalilova/authService/config"
	kafka "github.com/mirjalilova/authService/consumer"

	pb "github.com/mirjalilova/authService/genproto/auth"
	"github.com/mirjalilova/authService/service"
	"github.com/mirjalilova/authService/storage/postgres"

	"path/filepath"
	"runtime"

	"google.golang.org/grpc"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func main() {
	cnf := cf.Load()
	db, err := postgres.NewPostgresStorage(cnf)
	if err != nil {
		slog.Error("can't connect to db: %v", err)
		return
	}
	defer db.Db.Close()

	userService := service.NewUserService(db)

	brokers := []string{"localhost:9092"}

	kcm := kafka.NewKafkaConsumerManager()

	if err := kcm.RegisterConsumer(brokers, "upd-evuseral", "eval", kafka.UserEditProfileHandler(userService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'upd-user' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}
	
	listener, err := net.Listen("tcp", cnf.USER_PORT)
	if err != nil {
		slog.Error("can't listen: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, service.NewUserService(db))

	slog.Info("server started port", cnf.USER_PORT)
	if err := s.Serve(listener); err != nil {
		slog.Error("can't serve: %v", err)
		return
	}
}
