package model

import "time"

// Make separate request and response model and keep the fields accordingly.

type Todos struct {
	Id              int    `json:"id" db:"id, primary key"`
	TodoName        string `json:"todoName" db:"todoname"`
	TodoDescription string `json:"todoDescription" db:"tododescription"`

	/*
		JSON format should be camel case and db format should be snake case
	*/

	IsCompleted bool      `json:"iscompleted" db:"iscompleted"`
	UserId      int       `json:"user_id" db:"userId"` //Foregin key for user who owns this todo
	CreatedAt   time.Time `json:"created_at" db:"createdat"`
	UpdateAt    time.Time `json:"updateAt" db:"updateat"`
}

type DeleteTodos struct {
	Id int `json:"id" db:"id, primary key"`
}

//type UpdateTodos struct {
//	Id              string    `json:"id" db:"id, primary key"`
//	TodoName        string    `json:"todoName" db:"todoname"`
//	TodoDescription string    `json:"todoDescription" db:"tododescription"`
//	IsCompleted     bool      `json:"iscompleted" db:"iscompleted"`
//	UpdateAt        time.Time `json:"updateAt" db:"updateat"`
//}

//type user struct {
//	username
//}
