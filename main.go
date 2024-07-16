package main

import (
	"net"

	"golang.org/x/exp/slog"

	cf "github.com/mirjalilova/authService/config"

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
	config := cf.Load()
	db, err := postgres.NewPostgresStorage(config)
	if err != nil {
		slog.Error("can't connect to db: %v", err)
		return
	}
	defer db.Db.Close()

	listener, err := net.Listen("tcp", config.AUTH_PORT)
	if err != nil {
		slog.Error("can't listen: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, service.NewAuthService(db))
	pb.RegisterUserServiceServer(s, service.NewUserService(db))

	slog.Info("server started port", config.COMPANY_PORT)
	if err := s.Serve(listener); err != nil {
		slog.Error("can't serve: %v", err)
		return
	}
}
