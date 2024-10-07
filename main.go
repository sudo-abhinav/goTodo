package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/sudo-abhinav/go-todo/handler"
	"net/http"
)

func main() {

	/*
		1) try to keep main.go in a cmd folder
		2) group public routes that do not need authentication like register, login, healthcheck
		3) group protected routes that needs authentication
		4) Need to add middleware for authentication
		5) Create a separate DB call
		6) Create a separate routes file in srever.go and placed the file in server folder
	*/

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// you can directly respond here only like sever is running...

	r.Get("/healthcheck", handler.HealthChecker)
	r.Route("/api", func(r chi.Router) {

		// routes name should be relevant so that we can identify the task

		r.Get("/data", handler.GetAllTodo) // get All todo
		r.Get("/databyid/{id}", handler.GetTodoById)
		r.Post("/createTodo", handler.CreateTodo)
		r.Delete(`/deleteById/{id}`, handler.DeleteTodoById)
		r.Put("/updatetodo", handler.UpdateTodo)

		//	user Registration
		r.Post("/reg", handler.UserRegstration)
		r.Post("/login", handler.UserLogin)
	})

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
