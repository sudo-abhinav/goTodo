package response

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		return
	}
}
