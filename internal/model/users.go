package model

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string
}

// Пример "базы данных" (замените на реальную БД)
var users = map[string]User{
	"admin": {
		ID:           1,
		Username:     "admin",
		PasswordHash: "$2a$10$3Y8V7y2Z5q6Q5h5d5v5Z3e",
		Role:         "admin",
	},
	"user": {
		ID:           2,
		Username:     "user",
		PasswordHash: "",
		Role:         "user",
	},
}

func GetUser(username string) (User, error) {
	user, exists := users[username]
	if !exists {
		return User{}, fmt.Errorf("User not found")
	}
	return user, nil
}

func (u *User) CheckPassword(password string) bool {
	hash, err := bcrypt.GenerateFromPassword([]byte("admin"), 14)
	if err != nil {
		log.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
