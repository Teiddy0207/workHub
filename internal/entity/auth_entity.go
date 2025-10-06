package entity

import "time"

type User struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Username  string    `json:"username"`
    Password  string    `json:"-"`
    CreatedAt time.Time `json:"created_at"`
}