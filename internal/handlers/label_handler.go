package handlers

import (
	"github.com/milwad-dev/do-it/internal/models"
	"github.com/milwad-dev/do-it/internal/utils"
	"net/http"
)

// StoreLabel => store new label and return json response
func (db *DBHandler) StoreLabel(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	r.ParseForm()
	// TODO: ADD VALIDATION
	// Read request body
	var label *models.Label
	utils.ReadRequestBody(w, r, label)
	userId := 1 // TODO FIX THIS

	// Exec query (id, created_at and updated_at filled automatically)
	query := "INSERT INTO labels (title, color, user_id) VALUES (?, ?, ?)"
	_, err := db.Exec(query, &label.Title, &label.Color, userId)
	if err != nil {
		panic(err)
	}

	data := make(map[string]string)
	data["message"] = "The label store successfully."

	utils.JsonResponse(w, data, 200)
}
