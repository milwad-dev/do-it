package handlers

import (
	"encoding/json"
	"github.com/milwad-dev/do-it/internal/models"
	"github.com/milwad-dev/do-it/internal/utils"
	"net/http"
)

// GetLatestLabels => Get the latest labels and return json response
// @Summary Get Labels
// @Description Get the latest labels
// @Produce json
// @Success 200 {object} []models.Label
// @Router /api/labels [get]
func (db *DBHandler) GetLatestLabels(w http.ResponseWriter, r *http.Request) {
	var labels []models.Label

	query := `
    SELECT l.id, l.title, l.color, l.created_at, l.updated_at, l.user_id, 
           u.id, u.name, COALESCE(u.email, ''), COALESCE(u.phone, ''), u.created_at 
    FROM labels l
    JOIN users u ON l.user_id = u.id
    ORDER BY l.created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var label models.Label
		var user models.User

		err := rows.Scan(
			&label.ID, &label.Title, &label.Color, &label.CreatedAt, &label.UpdatedAt, &label.UserId,
			&user.ID, &user.Name, &user.Email, &user.Phone, &user.CreatedAt,
		)
		if err != nil {
			panic(err)
		}

		label.User = user
		labels = append(labels, label)
	}

	// TODO: Fix the format of json
	utils.JsonResponse(w, labels, 200)
}

// StoreLabel => store new label and return json response
// @Summary Store Label
// @Description store new label
// @Produce json
// @Param title body string true "The title of the label"
// @Param color body string true "The color of the label"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/labels [post]
func (db *DBHandler) StoreLabel(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	r.ParseForm()

	// TODO: ADD VALIDATION

	// Read request body
	var label models.Label

	userId := 9 // TODO FIX THIS

	// Decode JSON request body into `labels`
	err := json.NewDecoder(r.Body).Decode(&label)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Exec query (id, created_at and updated_at filled automatically)
	query := "INSERT INTO labels (title, color, user_id) VALUES (?, ?, ?)"
	_, err = db.Exec(query, &label.Title, &label.Color, userId)
	if err != nil {
		panic(err)
	}

	data := make(map[string]string)
	data["message"] = "The label store successfully."

	utils.JsonResponse(w, data, 200)
}
