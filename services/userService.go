package services

import (
	"fmt"
	"github.com/sudo-abhinav/go-todo/Database"
	"github.com/sudo-abhinav/go-todo/model"
	"github.com/sudo-abhinav/go-todo/utils/encryption"
	"time"
)

func CreateUserInDB(user model.UserReg) error {
	//here i am creating hash password for security reason
	hashPWD, err := encryption.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Prepare the SQL query
	queryString := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`

	res, err := Database.DBconn.Exec(queryString, user.UserName, user.Email, hashPWD)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("no User created")
	}

	return nil

}

func CreateTodoInDb(todos model.Todos) error {
	QueryString := "insert into usertodo (todoname, tododescription  ) VALUES ($1,$2 ) "

	res, err := Database.DBconn.Exec(QueryString, todos.TodoName, todos.TodoDescription)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("no todo crated")
	}

	return nil

}

func UpdateTodoInDB(todos model.Todos) error {
	query := "UPDATE usertodo SET todoname = $1, tododescription = $2, updateat = $3 WHERE id = $4"
	res, err := Database.DBconn.Exec(query, todos.TodoName, todos.TodoDescription, time.Now(), todos.Id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("no todo found with id %d", todos.Id)
	}

	return nil
}

func DeleteTodoInDB(param model.DeleteTodos) error {
	res, err := Database.DBconn.Exec("DELETE FROM usertodo WHERE id=$1", param.Id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("No Todo Found with id %d", param.Id)
	}
	return nil
}
