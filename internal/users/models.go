package users

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const hashCost = bcrypt.DefaultCost

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	// Here we will have a hashed password
	Password string `json:"-"`

	IsAdmin bool `json:"isAdmin"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(name, password string) (User, error) {
	if len(name) == 0  || len(password) == 0{
		return User{}, errors.New("users: name and password cant be empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return User{}, err
	}
	user := User{
		Name:      name,
		Password:  string(hashedPassword),
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}
	return user, nil
}

func (u *User)ComparePasswords(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func HashePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), hashCost)
}