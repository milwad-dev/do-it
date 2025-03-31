package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JsonResponse => set json header and return json response
func JsonResponse(w http.ResponseWriter, data any, statusCode int) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Set status code of response
	w.WriteHeader(statusCode)

	// Set status from status code
	status := "Success"
	if statusCode >= 400 && statusCode < 500 {
		status = "Client Error"

		// TODO: Add log
	} else if statusCode >= 500 && statusCode < 600 {
		status = "Server Error"

		// TODO: Add log
	}

	// Attach status to data if data is map
	response := make(map[string]interface{})
	switch v := data.(type) {
	case map[string]interface{}:
		response = v
	default:
		fmt.Println(v)
		response["data"] = v
	}
	response["status"] = status

	// Return json response
	json.NewEncoder(w).Encode(response)
}
