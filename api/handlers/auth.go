package handlers

import (
	"github.com/mirjalilova/authService/genproto/auth"
	"context"
	t "github.com/mirjalilova/authService/api/token"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


// RegisterUser handles user registration
// @Summary Register a new user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body auth.RegisterReq true "Register User Request"
// @Success 200 {object} auth.RegisterUserRes
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /register [post]
func (h *Handlers) RegisterUser(c *gin.Context) {
	var req auth.RegisterReq
	if err := c.BindJSON(&req); err != nil {
		log.Printf("failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	password, err := t.HashPassword(req.Password)
	if err!= nil {
        log.Printf("failed to hash password: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }

	req.Password = password

	_, err = h.Auth.Register(context.Background(), &req)
	if err != nil {
		log.Printf("failed to register user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// LoginUser handles user login
// @Summary Login a user
// @Description Login a user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body auth.LoginReq true "Login Request"
// @Success 200 {string} auth.User
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /login [post]
func (h *Handlers) LoginUser(c *gin.Context) {
	var req auth.LoginReq
	if err := c.BindJSON(&req); err != nil {
		log.Printf("failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	res, err := h.Auth.Login(context.Background(), &req)
	if err != nil {
		log.Printf("failed to login user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	token := t.GenerateJWTToken(res)

	c.JSON(http.StatusOK, token)
}

func (h *Handlers) ForgotPassword(c *gin.Context) {
	var req auth.GetByEmail
    if err := c.BindJSON(&req); err != nil {
        log.Printf("failed to bind JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    _, err := h.Auth.ForgotPassword(context.Background(), &req)
    if err != nil {
        log.Printf("failed to send password reset email: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent successfully"})
}

func (h *Handlers) ResetPassword(c *gin.Context) {
	var req auth.ResetPassReq
    if err := c.BindJSON(&req); err != nil {
        log.Printf("failed to bind JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    _, err := h.Auth.ResetPassword(context.Background(), &req)
    if err != nil {
        log.Printf("failed to reset password: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}