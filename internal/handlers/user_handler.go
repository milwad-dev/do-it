package handlers

import (
	"github.com/milwad-dev/do-it/internal/models"
	"github.com/milwad-dev/do-it/internal/utils"
	"log"
	"net/http"
)

// GetLatestUsers => get the latest users and return json response
// @Summary Get Users
// @Description Get the latest users
// @Produce json
// @Success 200 {object} []models.User
// @Router /api/users [get]
func (db *DBHandler) GetLatestUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	sql := "SELECT id, name, COALESCE(email, ''), COALESCE(phone, ''), created_at FROM USERS ORDER BY created_at DESC"
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.CreatedAt)
		if err != nil {
			panic(err.Error())
		}

		users = append(users, user)
	}

	// TODO: Fix the format of json
	utils.JsonResponse(w, users, 200)
}
