package handlers

import (
	"github.com/milwad-dev/do-it/internal/services"
	"github.com/milwad-dev/do-it/internal/utils"
	"log"
	"net/http"
	"net/mail"
)

func (db *DBHandler) RegisterAuth(w http.ResponseWriter, r *http.Request) {
	name := r.Form.Get("name")
	username := r.Form.Get("username") // Email or Phone
	password := r.Form.Get("password")

	usernameField := detectEmailOrPhone(username)

	rows, err := db.Query("SELECT count(*) FROM users WHERE ? = ?", usernameField, username)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var count int
	var data map[string]any

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	// Check user is exists in db
	if count == 1 {
		data["message"] = "The user already exists."

		utils.JsonResponse(w, data, 302)
	}

	// Hash password
	hashPassword, _ := services.HashPassword(password)

	// Create user
	result, err := db.Exec("INSERT INTO users (name, ?, password) VALUES (?, ?, ?)", usernameField, name, username, hashPassword)
	if err != nil {
		data["message"] = "Problem on creating user."

		utils.JsonResponse(w, data, 302)
	}

	// Create token
	userId, _ := result.LastInsertId()
	token, err := services.GenerateToken(uint(userId))
	if err != nil {
		data["message"] = "Problem on generating token."

		utils.JsonResponse(w, data, 302)
	}

	data["message"] = "Register completed."
	data["token"] = token

	utils.JsonResponse(w, data, 200)
}

func (db *DBHandler) LoginAuth(w http.ResponseWriter, r *http.Request) {
	//name := r.Form.Get("name")
	//username := r.Form.Get("username") // Email or Phone
	//password := r.Form.Get("password")
	//
	//usernameField := detectEmailOrPhone(username)
	//
	//rows, err := db.Query("SELECT count(*) FROM users WHERE ? = ?", usernameField, username)
	//if err != nil {
	//	panic(err)
	//}
	//defer rows.Close()
	//
	//var count int
	//
	//for rows.Next() {
	//	if err := rows.Scan(&count); err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//if count == 0 {
	//	data := map[string]any{
	//		"message": "The user not found",
	//	}
	//	utils.JsonResponse(w, data, 302)
	//}
	//
	//data := map[string]any{
	//	"message": "Register completed",
	//}
	//utils.JsonResponse(w, data, 302)
}

// detectEmailOrPhone => Detect the username is email or phone
func detectEmailOrPhone(username string) string {
	_, err := mail.ParseAddress(username)
	if err != nil {
		return "phone"
	}

	return "email"
}
