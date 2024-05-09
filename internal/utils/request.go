package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

// ReadRequestBody => read and parse request body
func ReadRequestBody(w http.ResponseWriter, r *http.Request, model interface{}) interface{} {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return nil
	}

	// Parse request body
	err = json.Unmarshal(body, &model)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return nil
	}

	return model
}
