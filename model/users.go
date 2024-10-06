package model

import (
	"time"
)

type UserReg struct {
	Id        int       `json:"id" db:"id , Primary key"`
	UserName  string    `json:"username"  db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"createdat"`
	UpdateAt  time.Time `json:"update_at" db:"updateat"`
}
