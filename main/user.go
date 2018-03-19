package main

import (
	"time"
)

type User struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Birthdate time.Time `json:"birthDate"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
}
