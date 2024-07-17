package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/mirjalilova/authService/genproto/auth"
	"github.com/mirjalilova/authService/service"
)

func UserEditProfileHandler(u *service.UserService) func(message []byte) {
	return func(message []byte) {
		var user pb.UserRes
		if err := json.Unmarshal(message, &user); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		_, err := u.EditProfile(context.Background(), &user)
		if err != nil {
			log.Printf("failed to edit user via Kafka: %v", err)
			return
		}
		log.Printf("User updated")
	}
}