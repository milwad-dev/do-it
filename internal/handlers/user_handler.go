package handlers

import (
	"github.com/milwad-dev/do-it/internal/models"
	"github.com/milwad-dev/do-it/internal/utils"
	"log"
	"net/http"
)

// GetLatestUsers => get the latest users and return json response
func (db *DBHandler) GetLatestUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	sql := "SELECT id, name, email, phone, created_at FROM USERS ORDER BY created_at DESC"
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

	utils.JsonResponse(w, users, 200)
}
