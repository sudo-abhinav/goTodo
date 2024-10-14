package model

// use validate for input validation and make two models separately for request and response

// 1.todo do not use primary key and it is use serial in db auto which is auto genrated
type UserReg struct {
	Id       string `json:"id" db:"id "`
	UserName string `json:"username"  db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Login struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
