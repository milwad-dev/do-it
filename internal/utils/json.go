package utils

import (
	"encoding/json"
	"net/http"
)

// JsonResponse => set json header and return json response
func JsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
