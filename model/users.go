package model

import "time"

// Make separate request and response model and keep the fields accordingly.

type UserReg struct {
	Id       int    `json:"id" db:"id , Primary key"`
	UserName string `json:"username"  db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`

	/*
		JSON format should be camel case and db format should be snake case
	*/

	CreatedAt time.Time `json:"created_at" db:"createdat"`
	UpdateAt  time.Time `json:"update_at" db:"updateat"`
}
