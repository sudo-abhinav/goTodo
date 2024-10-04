package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/sudo-abhinav/go-todo/handler"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/healthcheck", handler.HealthChecker)
	r.Route("/api", func(r chi.Router) {

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
