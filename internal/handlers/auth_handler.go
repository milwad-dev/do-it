package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/milwad-dev/do-it/internal/services"
	"github.com/milwad-dev/do-it/internal/utils"
	"log"
	"net/http"
	"net/mail"
)

type user struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (db *DBHandler) RegisterAuth(w http.ResponseWriter, r *http.Request) {
	var user user

	data := make(map[string]any)
	// Parse body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, 400)
		return
	}

	// User info
	name := user.Name
	username := user.Username // Email or phone
	password := user.Password

	// Detect username
	usernameField := detectEmailOrPhone(username)

	// Check user not exists in db
	queryExist := fmt.Sprintf("SELECT count(*) FROM users WHERE %s = ?", usernameField)
	rows, err := db.Query(queryExist, username)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	// Check user is exists in db
	if count == 1 {
		data["message"] = "The user already exists."

		utils.JsonResponse(w, data, 302)
		return
	}

	// Hash password
	hashPassword, _ := services.HashPassword(password)

	// Create user
	query := fmt.Sprintf("INSERT INTO users (name, %s, password) VALUES (?, ?, ?)", usernameField)
	result, errInsert := db.Exec(query, name, username, hashPassword)
	if errInsert != nil {
		data["message"] = fmt.Sprintf("Problem on creating user: %s", errInsert.Error())

		utils.JsonResponse(w, data, 302)
		return
	}

	// Create token
	fmt.Println(name, username, password, usernameField, count)
	userId, _ := result.LastInsertId()
	token, errToken := services.GenerateToken(uint(userId))
	if errToken != nil {
		data["message"] = "Problem on generating token."

		utils.JsonResponse(w, data, 302)
		return
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
