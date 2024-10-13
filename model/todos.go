package model

type Todos struct {
	Id              string `json:"id" db:"id"`
	TodoName        string `json:"todoName" db:"todoname"`
	TodoDescription string `json:"todoDescription" db:"tododescription"`
	IsCompleted     bool   `json:"is_completed" db:"is_completed"`
	UserId          string `json:"userid" db:"userid"` //Foregin key for user who owns this todo
}

type DeleteTodos struct {
	Id string `json:"id" db:"id"`
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
