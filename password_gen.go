package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	var psw string
	fmt.Scan(&psw)
	s, _ := bcrypt.GenerateFromPassword([]byte(psw), bcrypt.DefaultCost)
	fmt.Println(string(s))
}