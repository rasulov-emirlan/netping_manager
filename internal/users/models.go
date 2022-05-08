package users

import (
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`

	IsAdmin bool `json:"isAdmin"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(name, password string) (User, error) {
	user := User{
		Name:      name,
		Password:  password,
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}
	return user, nil
}
