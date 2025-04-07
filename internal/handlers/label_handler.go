package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/milwad-dev/do-it/internal/models"
	"github.com/milwad-dev/do-it/internal/repositories"
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

	// Get user id from context
	userId := repositories.GetUserIdFromContext(r)

	// Create data
	data := make(map[string]interface{})

	query := `
    SELECT l.id, l.title, l.color, l.created_at, l.updated_at, l.user_id, 
           u.id, u.name, COALESCE(u.email, ''), COALESCE(u.phone, ''), u.created_at 
    FROM labels l
    JOIN users u ON l.user_id = u.id
    ORDER BY l.created_at DESC
    WHERE user_id = ?`
	rows, err := db.Query(query, userId)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusBadRequest)
		return
	}

	for rows.Next() {
		var label models.Label
		var user models.User

		err := rows.Scan(
			&label.ID, &label.Title, &label.Color, &label.CreatedAt, &label.UpdatedAt, &label.UserId,
			&user.ID, &user.Name, &user.Email, &user.Phone, &user.CreatedAt,
		)
		if err != nil {
			data["message"] = err.Error()

			utils.JsonResponse(w, data, http.StatusBadRequest)
			return
		}

		label.User = user
		labels = append(labels, label)
	}

	data["data"] = labels

	utils.JsonResponse(w, data, http.StatusOK)
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

	data := make(map[string]string)

	// Read request body
	var label models.Label

	// Get user id from context
	userId := repositories.GetUserIdFromContext(r)

	// Decode JSON request body into `labels`
	err := json.NewDecoder(r.Body).Decode(&label)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Create a new validator instance
	validate := validator.New()

	// Validate the User struct
	err = validate.Struct(label)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
		return
	}

	// Exec query (id, created_at and updated_at filled automatically)
	query := "INSERT INTO labels (title, color, user_id) VALUES (?, ?, ?)"
	_, err = db.Exec(query, &label.Title, &label.Color, userId)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	data["message"] = "The label store successfully."

	utils.JsonResponse(w, data, http.StatusOK)
}

// DeleteLabel => delete the label by id and return json response
// @Summary Delete Label
// @Description delete label by id
// @Produce json
// @Param id url integer true "The id of the label"
// @Success 200 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/labels/{id} [delete]
func (db *DBHandler) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	// Get label id from url
	labelId := chi.URLParam(r, "id")

	// Get user id from context
	userId := repositories.GetUserIdFromContext(r)

	// Create data
	data := make(map[string]string)

	queryExist := "SELECT count(*) FROM labels WHERE id = ? AND user_id = ?"
	rows, err := db.Query(queryExist, labelId, userId)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
		return
	}

	var count int

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			data["message"] = err.Error()

			utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
			return
		}
	}

	// If the label is not exists, we return error response
	if count == 0 {
		data["message"] = "The label is not exists."

		utils.JsonResponse(w, data, http.StatusNotFound)
		return
	}

	// If label exists, we delete it
	query := "DELETE FROM labels WHERE id = ? AND user_id = ?"
	_, err = db.Exec(query, labelId, userId)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
		return
	}

	data["message"] = "The label deleted successfully."

	utils.JsonResponse(w, data, http.StatusOK)
}
