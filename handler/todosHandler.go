package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sudo-abhinav/go-todo/Database"
	"github.com/sudo-abhinav/go-todo/Database/dbHelper"
	"github.com/sudo-abhinav/go-todo/middlewares"
	"github.com/sudo-abhinav/go-todo/model"
	"github.com/sudo-abhinav/go-todo/utils/encryption"
	"github.com/sudo-abhinav/go-todo/utils/response"
	_ "github.com/sudo-abhinav/go-todo/utils/response"
	"log"
	"net/http"
	"strconv"
)

type Claim struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func UserRegistration(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	//2.todo :  use validator

	if r.Body == nil {
		err := json.NewEncoder(w).Encode("please send some data")
		if err != nil {
			return
		}
	}
	var data model.UserReg
	userCtx := middlewares.UserContext(r)
	data.Id = userCtx.UserID

	//err := json.NewDecoder(r.Body).Decode(&data)
	if parseErr := response.ParseBody(r.Body, &data); parseErr != nil {
		response.RespondJSON(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if data.UserName == "" && data.Email == "" && data.Password == "" {
		response.RespondJSON(w, http.StatusBadRequest, "username, email, and password are required")
		return
	}
	if len(data.Password) <= 6 {
		response.RespondJSON(w, http.StatusBadRequest, "password length greater than 6 char")
		return
	}
	exists, existsErr := dbHelper.IsUserExists(data.Email)
	if err := dbHelper.CreateUserInDB(data.UserName, data.Email, data.Password); err != nil {
		response.RespondJSON(w, http.StatusInternalServerError, "error creating user")
		return
	}
	if existsErr != nil {
		response.RespondWithError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}

	if exists {
		response.RespondJSON(w, http.StatusBadRequest, "user already exists")
		return
	}
	response.RespondJSON(w, http.StatusCreated, "user Registered..")

}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Body)

	if r.Body == nil {
		//response.RespondJSON(w, http.StatusBadRequest, "please send some data")
		err := json.NewEncoder(w).Encode("please send some data")
		if err != nil {
			return
		}
	}

	var loginData model.Login
	if parseErr := response.ParseBody(r.Body, &loginData); parseErr != nil {
		response.RespondJSON(w, http.StatusBadRequest, "failed to parse request body")
		return
	}

	userID, username, userErr := dbHelper.GetUser(loginData.Email, loginData.Password)

	if userErr != nil {
		response.RespondJSON(w, http.StatusInternalServerError, "unauthorised access")
	}

	if userID == "" || username == "" {
		response.RespondWithError(w, http.StatusNotFound, nil, "user not found")
		return
	}

	SessionTicket, SessErr := dbHelper.CreateUserSession(userID)
	if SessErr != nil {
		response.RespondWithError(w, http.StatusInternalServerError, SessErr, "failed to create user session")
		return
	}

	tokenString, tokenErr := encryption.GenerateJWT(userID, loginData.Email, SessionTicket)
	if tokenErr != nil {
		response.RespondWithError(w, http.StatusInternalServerError, tokenErr, "failed to generate token")
		return
	}
	response.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{"user logged in successfully", tokenString})

	//TODO :-  am trying to add token in cookie
	//http.SetCookie(w, &http.Cookie{
	//	Name:    "Token",
	//	Value:   tokenString,
	//	Expires: time.Now().Add(time.Minute * 2),
	//})
	//response.RespondJSON(w, http.StatusCreated, "loggedIn")

}

func GetAllTodo(w http.ResponseWriter, r http.Request) {
	//w.Header().Set("content-Type", "application/json")
	//3.todo define [] using make keyword as possible because it creates by default null
	var posts []model.Todos
	//defer Database.DbConnectionClose()

	//4.todo use select function
	rows, err := Database.DBconn.Query("select usertodo.id , usertodo.todoname , usertodo.tododescription from usertodo")
	if err != nil {
		return
	}
	//log.Print(rows)
	for rows.Next() {
		//data := model.Todos{}
		var data model.Todos
		err := rows.Scan(&data.Id, &data.TodoName, &data.TodoDescription, &data.IsCompleted)
		if err != nil {
			fmt.Println("error", err)

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
	QueryString := `select id,todoname , tododescription , is_completed FROM usertodo WHERE id = $1`
	err := Database.DBconn.Get(&todo, QueryString, Param)
	if err != nil {
		//return http.StatusBadRequest
		w.Write([]byte("no data found"))
		w.WriteHeader(404)
		return

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

	var loginData model.Login
	if parseErr := response.ParseBody(r.Body, &loginData); parseErr != nil {
		response.RespondJSON(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	var todos model.Todos
	if err := json.NewDecoder(r.Body).Decode(&todos); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err, "invalid request payload")
		return
	}
	if err := dbHelper.CreateTodoInDb(todos); err != nil {
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
