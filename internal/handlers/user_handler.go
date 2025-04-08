package handlers

import (
	"github.com/milwad-dev/do-it/internal/models"
	"github.com/milwad-dev/do-it/internal/utils"
	"net/http"
)

// GetLatestUsers => get the latest users and return json response
// @Summary Get Users
// @Description Get the latest users
// @Produce json
// @Success 200 {object} []models.User
// @Router /users [get]
func (db *DBHandler) GetLatestUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	data := make(map[string]interface{})

	sql := "SELECT id, name, COALESCE(email, ''), COALESCE(phone, ''), created_at FROM USERS ORDER BY created_at DESC"
	rows, err := db.Query(sql)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusBadRequest)
		return
	}

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.CreatedAt)
		if err != nil {
			data["message"] = err.Error()

			utils.JsonResponse(w, data, http.StatusInternalServerError)
			return
		}

		users = append(users, user)
	}

	data["data"] = users

	utils.JsonResponse(w, data, http.StatusOK)
}
