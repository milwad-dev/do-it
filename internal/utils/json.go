package utils

import (
	"encoding/json"
	"net/http"
)

// JsonResponse => set json header and return json response
func JsonResponse(w http.ResponseWriter, data any, statusCode int) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Set status code of response
	w.WriteHeader(statusCode)

	// Return json response
	json.NewEncoder(w).Encode(data)
}
