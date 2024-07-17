package storage

import pb "github.com/mirjalilova/authService/genproto/auth"

type StorageI interface {
	Auth() AuthI
	User() UserI
}

type AuthI interface {
	Register(*pb.RegisterReq) (*pb.Void, error)
    Login(*pb.LoginReq) (*pb.User, error)
    ForgotPassword(*pb.GetByEmail) (*pb.Void, error)
    ResetPassword(*pb.ResetPassReq) (*pb.Void, error)
    RefreshToken(*pb.RefToken) (*pb.Void, error)
}

type UserI interface {
	GetProfile(*pb.GetById) (*pb.UserRes, error)
    EditProfile(*pb.UserRes) (*pb.UserRes, error)
    ChangePassword(*pb.ChangePasswordReq) (*pb.Void, error)
    GetSetting(*pb.GetById) (*pb.Setting, error)
    EditSetting(*pb.SettingReq) (*pb.Void, error)
    DeleteUser(*pb.GetById) (*pb.Void, error)
}