package handler

import (
	"fmt"
	"net/http"
)

// No need to create a function fot this

func HealthChecker(w http.ResponseWriter, r *http.Request) {

	// Send the response in json format
	// Make separate function for Encode JSON and Respond JSON

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
