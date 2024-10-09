package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sudo-abhinav/go-todo/Database"
	"github.com/sudo-abhinav/go-todo/model"
	"github.com/sudo-abhinav/go-todo/utils/encryption"
	"time"
)

var SecretKey = "secret-key-from-japan"

func CreateUserInDB(users model.UserReg) error {
	//here i am creating hash password for security reason
	hashPWD, err := encryption.HashPassword(users.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	queryString := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	res, err := Database.DBconn.Exec(queryString, users.UserName, users.Email, hashPWD)
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

func LoginUser(user model.UserReg) error {

	var hashPassword string
	var storedEmail string

	QueryString := `SELECT lower(email), password FROM users WHERE email = $1`
	err := Database.DBconn.QueryRow(QueryString, user.Email).Scan(&storedEmail, &hashPassword)
	if err != nil {
		return fmt.Errorf("error querying database: %w", err)
	}

	res := encryption.VerifyPassword(user.Password, hashPassword)
	if res == true && user.Email == storedEmail {
		//tokenString, err := auth.GenerateJWT(storedEmail)
		return nil
	}

	return nil
}

func IsUserExists(email string) (bool, error) {
	SQL := `SELECT email from users WHERE email = TRIM(LOWER($1))`
	var Useremail string
	err := Database.DBconn.Get(&Useremail, SQL, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, nil
}

func CreateTodoInDb(todos model.Todos) error {
	QueryString := "insert into usertodo (todoname, tododescription  ) VALUES ($1,$2 ) "
	res, err := Database.DBconn.Exec(QueryString, todos.TodoName, todos.TodoDescription)
	if err != nil {
		return err
	}

	//todo please ignore count
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

// 5.todo use archived _at column not hard delete the data
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
