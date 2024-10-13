package dbHelper

import (
	"fmt"
	"github.com/sudo-abhinav/go-todo/Database"
	"github.com/sudo-abhinav/go-todo/model"
	"time"
)

func CreateTodoInDb(todoname, tododescription, userID string) error {
	QueryString := "insert into usertodo (todoname, tododescription ,userid ) VALUES ($1,$2,$3) "
	_, err := Database.DBconn.Exec(QueryString, todoname, tododescription, userID)
	if err != nil {
		return err
	}
	return nil
	//todo please ignore count
	//count, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}
	//if count == 0 {
	//	return fmt.Errorf("no todo crated")
	//}
	//return nil
}

func UpdateTodoInDB(todos model.Todos) error {
	query := "UPDATE usertodo SET todoname = $1, tododescription = $2, update_at = $3  WHERE id = $4"
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
func DeleteTodoInDB(deleteTodo, userID string) error {
	query := `UPDATE usertodo SET
                    archived_at=now() 
                		WHERE  id= $1 AND 
                		    userid=$2 
                		  AND archived_at IS NULL`
	//res, err := Database.DBconn.Exec("DELETE FROM usertodo WHERE id=$1", param.Id)
	_, Err := Database.DBconn.Exec(query, deleteTodo, userID)
	if Err != nil {
		return Err
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

func GetCompleteTodos(userID string) ([]model.Todos, error) {
	query := `SELECT id, userid, todoname, tododescription, is_completed
			  FROM usertodo
			  WHERE userid = $1             
			    AND is_completed = true     
			    AND archived_at IS NULL`

	todos := make([]model.Todos, 0)
	Err := Database.DBconn.Select(&todos, query, userID)
	return todos, Err
}

func MarkCompleted(id, userID string) error {
	query := `UPDATE usertodo
              SET is_completed = true        
              WHERE id = $1                  
                AND userid = $2             
                AND archived_at IS NULL`

	_, Err := Database.DBconn.Exec(query, id, userID)
	if Err != nil {
		return Err
	}
	return nil
}

// GetAllTodos fetches all active todos for the specified user.
func GetAllTodos(UserID string) ([]model.Todos, error) {
	query := `SELECT id,  todoname, tododescription, is_completed,userid
			  FROM usertodo
			  WHERE userid = $1             
			    AND archived_at IS NULL`

	todos := make([]model.Todos, 0)
	fmt.Println(todos)
	Err := Database.DBconn.Select(&todos, query, UserID)
	return todos, Err
}

func DeleteAllTodos(userID string) (int, error) {
	query := `UPDATE usertodo
              SET archived_at = NOW()        
              WHERE userid = $1             
                AND archived_at IS NULL`

	_, delErr := Database.DBconn.Exec(query, userID)
	if delErr != nil {
		return 0, delErr
	}

	// If the query is successful, return a count of how many rows were expected to be affected
	// Here we assume that if the query was executed, it's for all relevant todos
	return 1, nil
}
