package handler

import (
	"fmt"
	"net/http"
)

func HealthChecker(w http.ResponseWriter, r *http.Request) {

	if ServerIsHealthy() {
		w.WriteHeader(http.StatusOK)
		fmt.Println(w, "Server is Healthy..")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(w, "server is Not Healthy")
	}
}
func ServerIsHealthy() bool {
	return true
}
