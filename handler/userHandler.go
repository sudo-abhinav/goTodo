package handler

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sudo-abhinav/go-todo/Database/dbHelper"
	"github.com/sudo-abhinav/go-todo/middlewares"
	"github.com/sudo-abhinav/go-todo/model"
	"github.com/sudo-abhinav/go-todo/utils/encryption"
	"github.com/sudo-abhinav/go-todo/utils/response"
	"net/http"
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
	//userCtx := middlewares.UserContext(r)
	//data.Id = userCtx.UserID

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

	userID, userEmail, userErr := dbHelper.GetUser(loginData.Email, loginData.Password)

	if userErr != nil {
		response.RespondJSON(w, http.StatusInternalServerError, "unauthorised access")
	}

	if userID == "" || userEmail == "" {
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

func DeActivateAccount(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID
	SessionId := userCtx.SessionID

	saveErr := dbHelper.DeleteUser(userID)
	if saveErr != nil {
		response.RespondWithError(w, http.StatusInternalServerError, saveErr, "failed to delete user account")
		return
	}

	saveErr = dbHelper.UserSessioneDelete(SessionId)
	if saveErr != nil {
		response.RespondWithError(w, http.StatusInternalServerError, saveErr, "failed to delete user session")
		return
	}

}
