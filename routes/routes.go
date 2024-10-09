package routes

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sudo-abhinav/go-todo/handler"
	"github.com/sudo-abhinav/go-todo/utils/response"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes() *Server {
	router := chi.NewRouter()

	router.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			response.RespondWithError(w, http.StatusOK, struct {
				Status string `json:"status"`
			}{Status: "server is running"})
		})
		r.Post("/reg", handler.UserRegstration)
		r.Post("/login", handler.UserLogin)
	})

	router.Group(func(r chi.Router) {
		//task:-here we use middlewre for authentication

		r.Get("/data", handler.GetAllTodo) // TODO :- this is only for testing purpose
		r.Get("/databyid/{id}", handler.GetTodoById)
		r.Post("/createTodo", handler.CreateTodo)
		r.Delete(`/deleteById/{id}`, handler.DeleteTodoById)
		r.Put("/updatetodo", handler.UpdateTodo)
	})
	return &Server{
		Router: router,
	}
}
func (server *Server) Run(port string) error {
	server.server = &http.Server{
		Addr:              port,
		Handler:           server.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return server.server.ListenAndServe()

}
func (server *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return server.server.Shutdown(ctx)
}
