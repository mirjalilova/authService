package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	pb "github.com/mirjalilova/authService/genproto/auth"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) GetProfile(req *pb.GetById) (*pb.UserRes, error) {
	res := &pb.UserRes{}

	var date string
	query := `SELECT id, username, email, full_name, date_of_birth FROM users WHERE id = $1`
	err := r.db.QueryRow(query, req.Id).
		Scan(
			&res.Id,
			&res.Username,
			&res.Email,
			&res.FullName,
			&date,
		)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	res.DateOfBirth = date[:10]

	return res, nil
}

func (r *UserRepo) EditProfile(req *pb.UserRes) (*pb.UserRes, error) {
	res := &pb.UserRes{}

	query := `UPDATE users SET updated_at = NOW()`

	var arg []interface{}
	var conditions []string

	if req.Username != "" && req.Username != "string" {
		arg = append(arg, req.Username)
		conditions = append(conditions, fmt.Sprintf("username = $%d", len(arg)))
	}

	if req.Email != "" && req.Email != "string" {
		arg = append(arg, req.Email)
		conditions = append(conditions, fmt.Sprintf("email = $%d", len(arg)))
	}

	if req.FullName != "" && req.FullName != "string" {
		arg = append(arg, req.FullName)
		conditions = append(conditions, fmt.Sprintf("full_name = $%d", len(arg)))
	}

	if req.DateOfBirth != "" && req.DateOfBirth != "string" {
		arg = append(arg, req.DateOfBirth)
		conditions = append(conditions, fmt.Sprintf("date_of_birth = $%d", len(arg)))
	}

	if len(conditions) > 0 {
		query += ", " + strings.Join(conditions, ", ")
	}

	query += fmt.Sprintf(" WHERE id = $%d", len(arg)+1)
	arg = append(arg, req.Id)

	_, err := r.db.Exec(query, arg...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *UserRepo) ChangePassword(req *pb.ChangePasswordReq) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE users SET updated_at = NOW(), password = $1 WHERE id = $2`
	_, err := r.db.Exec(query, req.NewPassword, req.Id)
	if err!= nil {
        return nil, err
    }

	return res, nil
}

func (r *UserRepo) GetSetting(req *pb.GetById) (*pb.Setting, error) {
	res := &pb.Setting{}

	query := `SELECT privacy_level, notification, language, theme FROM settings WHERE user_id = $1`
	err := r.db.QueryRow(query, req.Id).
		Scan(
			&res.PrivacyLevel,
            &res.Notification,
            &res.Language,
            &res.Theme,
		)
	if err == sql.ErrNoRows {
        return nil, err
    } else if err!= nil {
        return nil, err
    }

	return res, nil
}

func (r *UserRepo) EditSetting(req *pb.SettingReq) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE settings SET updated_at = NOW()`

	var arg []interface{}
	var conditions []string

	if req.PrivacyLevel!= "" && req.PrivacyLevel!= "string" {
        arg = append(arg, req.PrivacyLevel)
        conditions = append(conditions, fmt.Sprintf("privacy_level = $%d", len(arg)))
    }

	if req.Notification!= "" && req.Notification!= "string" {
        arg = append(arg, req.Notification)
        conditions = append(conditions, fmt.Sprintf("notification = $%d", len(arg)))
    }

	if req.Language!= "" && req.Language!= "string" {
        arg = append(arg, req.Language)
        conditions = append(conditions, fmt.Sprintf("language = $%d", len(arg)))
    }

	if req.Theme!= "" && req.Theme!= "string" {
        arg = append(arg, req.Theme)
        conditions = append(conditions, fmt.Sprintf("theme = $%d", len(arg)))
    }

	if len(conditions) > 0 {
        query += ", " + strings.Join(conditions, ", ")
    }

	query += fmt.Sprintf(" WHERE user_id = $%d", len(arg)+1)
	arg = append(arg, req.Id)
	_, err := r.db.Exec(query, arg...)
	if err!= nil {
        return nil, err
    }

	return res, nil
}

func (r *UserRepo) DeleteUser(req *pb.GetById) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE users SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1`
	_, err := r.db.Exec(query, req.Id)
	if err!= nil {
        return nil, err
    }

	return res, nil
}
