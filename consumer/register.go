package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/mirjalilova/authService/genproto/auth"
	"github.com/mirjalilova/authService/service"
)

func UserRegisterHandler(u *service.AuthService) func(message []byte) {
	return func(message []byte) {
		var user pb.RegisterReq
		if err := json.Unmarshal(message, &user); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		_, err := u.Register(context.Background(), &user)
		if err != nil {
			log.Printf("failed to register user via Kafka: %v", err)
			return
		}
		log.Printf("User registered")
	}
}