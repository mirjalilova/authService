package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	pb "github.com/mirjalilova/authService/genproto/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) Register(req *pb.RegisterReq) (*pb.Void, error) {
	res := &pb.Void{}

	query := `INSERT INTO users (username, email, password, full_name, date_of_birth) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, req.Username, req.Email, req.Password, req.FullName, req.DateOfBirth)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *AuthRepo) Login(req *pb.LoginReq) (*pb.User, error) {
	res := &pb.User{}

	var passwordHash string
	query := `SELECT id, username, email, password FROM users WHERE username = $1`
	err := r.db.QueryRow(query, req.Username).Scan(
		&res.Id,
		&res.Username,
		&res.Email,
		&passwordHash,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password for username: %s", req.Username)
	}

	return res, nil
}
func (r *AuthRepo) ForgotPassword(req *pb.GetByEmail) (*pb.Void, error) {
	res := &pb.Void{}

	query := `SELECT email FROM users WHERE email = $1`

	var email string
	err := r.db.QueryRow(query, req.Email).Scan(&email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s Email not found", req.Email)
		}
		return nil, err
	}

	return res, nil
}

func (r *AuthRepo) ResetPassword(req *pb.ResetPassReq) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE users SET password = $1, updated_at=now() WHERE email = $2`

	_, err := r.db.Exec(query, req.NewPassword, req.Email)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *AuthRepo) RefreshToken(req *pb.RefToken) (*pb.Void, error) {
	res := &pb.Void{}

    query := `INSERT INTO tokens (user_id, token) VALUES ($1, $2)`

    _, err := r.db.Exec(query, req.UserId, req.Token)
    if err != nil {
        return nil, err
    }

    return res, nil
}