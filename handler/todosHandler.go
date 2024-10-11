package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sudo-abhinav/go-todo/Database/dbHelper"
	"github.com/sudo-abhinav/go-todo/middlewares"
	"github.com/sudo-abhinav/go-todo/model"
	"github.com/sudo-abhinav/go-todo/utils/response"
	_ "github.com/sudo-abhinav/go-todo/utils/response"
	"net/http"
	"strconv"
)

//func GetAllTodo(w http.ResponseWriter, r http.Request) {
//	//w.Header().Set("content-Type", "application/json")
//	//3.todo define [] using make keyword as possible because it creates by default null
//	var posts []model.Todos
//	//defer Database.DbConnectionClose()
//
//	//4.todo use select function
//	rows, err := Database.DBconn.Query("select usertodo.id , usertodo.todoname , usertodo.tododescription from usertodo")
//	if err != nil {
//		return
//	}
//	//log.Print(rows)
//	for rows.Next() {
//		//data := model.Todos{}
//		var data model.Todos
//		err := rows.Scan(&data.Id, &data.TodoName, &data.TodoDescription, &data.IsCompleted)
//		if err != nil {
//			fmt.Println("error", err)
//
//		}
//		posts = append(posts, data)
//	}
//	log.Print(posts)
//	err = json.NewEncoder(w).Encode(posts)
//	if err != nil {
//		return
//	}
//}

func IncomoleteTodos(w http.ResponseWriter, r *http.Request) {
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

// completed services
func CreateTodo(w http.ResponseWriter, r *http.Request) {

	//log.Print("just print", r.Body)
	if r.Body == nil {
		json.NewEncoder(w).Encode("please send some data")

	}

	var todos model.Todos
	userCtx := middlewares.UserContext(r)
	todos.UserId = userCtx.UserID
	fmt.Println(todos.UserId, "this token came from context")
	if parseErr := response.ParseBody(r.Body, &todos); parseErr != nil {
		response.RespondJSON(w, http.StatusBadRequest, "failed to parse request body")
		return
	}

	fmt.Println(todos.TodoName, todos.TodoDescription, todos.UserId)

	if err := dbHelper.CreateTodoInDb(todos.TodoName, todos.TodoDescription, todos.UserId); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err, "error creating todo")
		return
	}
	response.RespondJSON(w, http.StatusCreated, "Todo Created..")

}

func DeleteTodoById(w http.ResponseWriter, r *http.Request) {

	param := chi.URLParam(r, "id")
	fmt.Println(param)

	id, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	deleteTodo := model.DeleteTodos{Id: id}

	err = dbHelper.DeleteTodoInDB(deleteTodo)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err, "error Deleting todo")
		return

	}

	response.RespondJSON(w, http.StatusOK, "Todo Update..")
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
	//json.NewEncoder(w).Encode("Todo Updated")
}
