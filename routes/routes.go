package routes

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	chi.Router
	server *http.Server
}

func SetupRoutes() {

}
