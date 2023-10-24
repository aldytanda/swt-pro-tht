// This file contains types that are used in the repository layer.
package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateColumn = errors.New("duplicate_column")
	ErrNotFound        = errors.New("not_found")
	ErrUnauthorized    = errors.New("unauthorized")
)

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type User struct {
	ID         int64
	Name       string
	Phone      string
	Password   string
	CountLogin int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Login struct {
	Phone    string
	Password string
}

type JWT struct {
	token string
}

func (u *User) hashPassword() error {
	password := []byte(u.Password)

	hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error while hashing password: %w", err)
	}

	u.Password = string(hashed)

	return nil
}

func (u *User) passwordMatch(inputPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword)); err == nil {
		return true
	}

	return false
}

func (u *User) generateToken(jwtSecretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(480 * time.Minute).Unix(),
		"authorized": true,
		"sub":        fmt.Sprint(u.ID),
		"user":       u.Name,
	})

	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("Error while generating token: %w", err)
	}

	return tokenString, nil
}
