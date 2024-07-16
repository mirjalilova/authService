package handlers

import (
	"github.com/mirjalilova/authService/service"
)

type Handlers struct {
	Auth *service.AuthService
}

func NewHandler(ah *service.AuthService) *Handlers {
	return &Handlers{
		Auth: ah,
	}
}
