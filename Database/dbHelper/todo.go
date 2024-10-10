package dbHelper

import (
	"fmt"
	"github.com/sudo-abhinav/go-todo/Database"
	"github.com/sudo-abhinav/go-todo/model"
	"time"
)

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
	query := "UPDATE usertodo SET todoname = $1, tododescription = $2, update_at = $3 WHERE id = $4"
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

func GetIncompleteTodos(userID string) ([]model.Todos, error) {
	query := `SELECT id, userid, todoname, tododescription, is_completed
			  FROM usertodo
			  WHERE userid = $1             
			    AND is_completed = false     
			    AND archived_at IS NULL`

	todos := make([]model.Todos, 0)
	Err := Database.DBconn.Select(&todos, query, userID)
	return todos, Err
}
