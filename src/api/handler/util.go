package handler

import (
	"encoding/json"
	"net/http"
)

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response, _ := json.Marshal(map[string][]map[string]string{
		"errors": {
			0: map[string]string{"message": message},
		},
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
