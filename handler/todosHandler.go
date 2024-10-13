package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/sudo-abhinav/go-todo/Database/dbHelper"
	"github.com/sudo-abhinav/go-todo/middlewares"
	"github.com/sudo-abhinav/go-todo/model"
	"github.com/sudo-abhinav/go-todo/utils/response"
	_ "github.com/sudo-abhinav/go-todo/utils/response"
	"net/http"
)

func InCompleteTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todos, err := dbHelper.GetIncompleteTodos(userID)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err, "error in getting todos")
	}

	if len(todos) == 0 {
		response.RespondJSON(w, http.StatusNotFound, "No todo found")
		return
	}

	response.RespondJSON(w, http.StatusCreated, todos)
}

// completed services
func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todos, Err := dbHelper.GetAllTodos(userID)
	fmt.Println(userID)
	fmt.Println(Err)
	if Err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, Err, "failed to get todos")
		return
	}
	if len(todos) == 0 {
		response.RespondWithError(w, http.StatusNotFound, Err, "no todo found")
		return
	}

	response.RespondJSON(w, http.StatusCreated, todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {

	var todos model.Todos

	if parseErr := response.ParseBody(r.Body, &todos); parseErr != nil {
		response.RespondWithError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	userCtx := middlewares.UserContext(r)
	UserId := userCtx.UserID
	fmt.Println(todos.UserId, "this token came from context")
	if parseErr := response.ParseBody(r.Body, &todos); parseErr != nil {
		response.RespondJSON(w, http.StatusBadRequest, "failed to parse request body")
		return
	}

	fmt.Println(todos.TodoName, todos.TodoDescription, todos.UserId)

	if err := dbHelper.CreateTodoInDb(todos.TodoName, todos.TodoDescription, UserId); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err, "error creating todo")
		return
	}
	response.RespondJSON(w, http.StatusCreated, "Todo Created..")

}

func DeleteTodoById(w http.ResponseWriter, r *http.Request) {

	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	param := chi.URLParam(r, "id")
	fmt.Println(param)

	if param == "" {
		response.RespondWithError(w, http.StatusBadRequest, nil, "Missing ID parameter")
		return
	}
	// Call the database helper to delete the todo item
	err := dbHelper.DeleteTodoInDB(param, userID)
	if err != nil {
		// Check if the error is due to the item not found or other reasons
		if errors.Is(err, sql.ErrNoRows) {
			response.RespondWithError(w, http.StatusNotFound, err, "Todo not found")
		} else {
			response.RespondWithError(w, http.StatusInternalServerError, err, "Error deleting todo")
		}
		return
	}
	response.RespondJSON(w, http.StatusOK, "Todo Deleted..")
	//json.NewEncoder(w).Encode("Todo Deleted")
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		response.RespondJSON(w, http.StatusBadRequest, "please send all data")
		return
	}
	var todos model.Todos
	if err := json.NewDecoder(r.Body).Decode(&todos); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err, "invalid request payload")
		return
	}
	if err := dbHelper.UpdateTodoInDB(todos); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err, "error updating todo")
		return
	}
	response.RespondJSON(w, http.StatusCreated, "Todo Update..")
	response.RespondJSON(w, http.StatusCreated, "Todo Update..")
	//json.NewEncoder(w).Encode("Todo Updated")
}

func GetComoleteTodo(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todos, err := dbHelper.GetCompleteTodos(userID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err, "error in Getting completed todo")
		return
	}
	if len(todos) == 0 {
		response.RespondWithError(w, http.StatusNotFound, err, "no todo found")
		return
	}
	response.RespondJSON(w, http.StatusCreated, todos)

}

func DeleteAll(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	res, err := dbHelper.DeleteAllTodos(userID)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err, "Error in Deleting All todo")
	}

	if res == 0 {
		response.RespondWithError(w, http.StatusNotFound, nil, "No todo found")
		return
	}

	response.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"all todos deleted successfully"})
}
