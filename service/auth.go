package service

import (
	"context"

	pb "github.com/mirjalilova/authService/genproto/auth"
	st "github.com/mirjalilova/authService/storage/postgres"
)

type AuthService struct {
	storage st.Storage
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(storage *st.Storage) *AuthService {
	return &AuthService{
		storage: *storage,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.Void, error) {
	res, err := s.storage.AuthS.Register(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginReq) (*pb.User, error) {
	res, err := s.storage.AuthS.Login(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}


func (s *AuthService) ForgotPassword(ctx context.Context, req *pb.GetByEmail) (*pb.Void, error) {
	res, err := s.storage.AuthS.ForgotPassword(req)
    if err != nil {
        return nil, err
    }
    return res, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *pb.ResetPassReq) (*pb.Void, error) {
	res, err := s.storage.AuthS.ResetPassword(req)
    if err != nil {
        return nil, err
    }
    return res, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefToken) (*pb.Void, error) {
	res, err := s.storage.AuthS.RefreshToken(req)
    if err != nil {
        return nil, err
    }
    return res, nil
}