package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgx.Conn
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	GetUser(id int) (*User, error)
	UpdatePassword(id int, currentPassword, newPassword string) error
}

func (m *UserModel) Insert(name, email, password string) error {
	stmt := `
    INSERT INTO users (name, email, hashed_password, created)
    VALUES($1, $2, $3, CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
  `
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec(context.Background(), stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = $1"
	err := m.DB.QueryRow(context.Background(), stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) GetUser(id int) (*User, error) {
	var user User
	stmt := `SELECT name, email, created FROM users WHERE id = $1`
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &user, ErrNoRecord
		} else {
			return &user, err
		}
	}
	return &user, nil
}

func (m *UserModel) UpdatePassword(id int, currentPassword, newPassword string) error {
	var hashedCurrentPassword []byte

	stmt := "SELECT hashed_password FROM users WHERE id = $1"
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&hashedCurrentPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedCurrentPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	stmt = `UPDATE users SET hashed_password=$1 WHERE id=$2`
	_, err = m.DB.Exec(context.Background(), stmt, string(hashedNewPassword), id)
	return err
}
