package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sudo-abhinav/go-todo/Database"
	"github.com/sudo-abhinav/go-todo/model"
	"github.com/sudo-abhinav/go-todo/services"
	"github.com/sudo-abhinav/go-todo/utils/response"
	_ "github.com/sudo-abhinav/go-todo/utils/response"
	"log"
	"net/http"
	"strconv"
)

/*
	Put user related handler in user.go and todo related in todo.go
*/

func UserRegstration(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

	// No need to print such kind of things try to log relevant things

	log.Printf("Received %s request for %s", r.Method, r.URL.Path)

	// Create a separate parseBody function in utils.go and call it whenever needed

	log.Print("just print ", r.Body)
	if r.Body == nil {
		err := json.NewEncoder(w).Encode("please send some data")
		if err != nil {
			return
		}
	}
	var data model.UserReg
	err := json.NewDecoder(r.Body).Decode(&data)
	fmt.Println(data)
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if data.UserName == "" || data.Email == "" || data.Password == "" {
		response.RespondWithError(w, http.StatusBadRequest, "username, email, and password are required")
		return
	}

	if err := services.CreateUserInDB(data); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	// Create a separate EncodeJSON and RespondJSON function in utils.go
	// Send status code also in response

	json.NewEncoder(w).Encode("user Created")

}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Body)

	if r.Body == nil {
		err := json.NewEncoder(w).Encode("please send some data")
		if err != nil {
			return
		}
	}
	var loginData model.UserReg
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "invalid payload request")
	}

	if err := services.LoginUser(loginData); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "error creating user")
	}
	json.NewEncoder(w).Encode("LoggedIN")

	/*
		1) You are not using JWT and session management
	*/
}

func GetAllTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	var posts []model.Todos
	//defer Database.DbConnectionClose()

	// Make the separate function for every query and put it in dbHelper folder

	rows, err := Database.DBconn.Query("select usertodo.id , usertodo.todoname , usertodo.tododescription , usertodo.iscompleted from usertodo")
	if err != nil {
		log.Fatal(err)
	}

	// Use select function and store all the results in a model once

	//log.Print(rows)
	for rows.Next() {
		//data := model.Todos{}
		var data model.Todos
		err := rows.Scan(&data.Id, &data.TodoName, &data.TodoDescription, &data.IsCompleted)
		if err != nil {
			fmt.Println("error", err)
			log.Fatal(err)
		}
		posts = append(posts, data)
	}
	log.Print(posts)
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		return
	}
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("content-Type", "application/json")
	var todo model.Todos
	Param := chi.URLParam(r, "id")
	//fmt.Println(Param)
	QueryString := `select id,todoname , tododescription , iscompleted FROM usertodo WHERE id = $1`
	err := Database.DBconn.Get(&todo, QueryString, Param)
	if err != nil {
		//return http.StatusBadRequest
		w.Write([]byte("no data found"))

		// status code should be 200 - StatusOK

		w.WriteHeader(404)
		return
		//log.Fatal("error in Query string : ", err)
	}
	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		return
	}

}

// completed services
func CreateTodo(w http.ResponseWriter, r *http.Request) {

	//log.Print("just print", r.Body)
	if r.Body == nil {
		json.NewEncoder(w).Encode("please send some data")
	}
	var todos model.Todos
	if err := json.NewDecoder(r.Body).Decode(&todos); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := services.CreateTodoInDb(todos); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "error creating todo")
		return
	}
	json.NewEncoder(w).Encode("Todo Created")

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

	err = services.DeleteTodoInDB(deleteTodo)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "error Deleting todo")
		return

	}
	json.NewEncoder(w).Encode("Todo Deleted")
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		response.RespondWithError(w, http.StatusBadRequest, "please send all data")
		return
	}
	var todos model.Todos
	if err := json.NewDecoder(r.Body).Decode(&todos); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := services.UpdateTodoInDB(todos); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "error updating todo")
		return
	}
	json.NewEncoder(w).Encode("Todo Updated")
}
