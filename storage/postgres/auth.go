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

	query := `INSERT INTO users (username, email, password, full_name, role, date_of_birth) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, req.Username, req.Email, req.Password, req.FullName, req.Role, req.DateOfBirth)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *AuthRepo) Login(req *pb.LoginReq) (*pb.User, error) {
	res := &pb.User{}

	var password string
	query := `SELECT id, username, email, password, role ROM users WHERE username = $1 AND password = $2`
	err := r.db.QueryRow(query, req.Username, req.Password).
		Scan(
			&res.Id,
			&res.Username,
			&res.Email,
			&password,
			&res.Role,
		)

	if req.Username != res.Username {
		return nil, fmt.Errorf("nivalid username: %s", req.Username)
	}
	err = bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password: %s", req.Username)
	}
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *AuthRepo) ForgotPassword(req *pb.GetByEmail) (*pb.Void, error) {
	res := &pb.Void{}

	query := `SELECT email FROM users WHERE email = $1`

	var email string
	err := r.db.QueryRow(query, req.Email).Scan(email)

	if email != req.Email {
		return nil, fmt.Errorf("%s Email not found", req.Email)
	}

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

	query := `UPDATE users SET password = $1 WHERE id = $2`

	_, err := r.db.Exec(query, req.NewPassword, req.Id)
	if err != nil {
		return nil, err
	}

	return res, nil
}
