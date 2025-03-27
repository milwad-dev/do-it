package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/milwad-dev/do-it/internal/services"
	"github.com/milwad-dev/do-it/internal/utils"
	"log"
	"net/http"
	"net/mail"
	"regexp"
)

// RegisterAuth => Register user and create token
// @Summary Register user
// @Description Create new user with token
// @Produce json
// @Param name body string true "The name of the user"
// @Param username body string true "The email or phone of the user"
// @Param password body string true "The password of the user"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/register [post]
func (db *DBHandler) RegisterAuth(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// TODO: ADD validation
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
	userId, _ := result.LastInsertId()
	token, errToken := services.GenerateToken(uint(userId))
	if errToken != nil {
		data["message"] = "Problem on generating token."

		utils.JsonResponse(w, data, 302)
		return
	}

	// Return success response
	data["message"] = "Register completed."
	data["token"] = token

	utils.JsonResponse(w, data, 200)
}

// LoginAuth => Check user credentials and create jwt token
// @Summary Login user
// @Description Check user credentials and login
// @Produce json
// @Param username body string true "The email or phone of the user"
// @Param password body string true "The password of the user"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/login [post]
func (db *DBHandler) LoginAuth(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	data := make(map[string]any)

	// Parse body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data["message"] = err.Error()
		utils.JsonResponse(w, data, 400)
		return
	}

	// Detect username type (email or phone)
	usernameField := detectEmailOrPhone(user.Username)

	// Get user info from database
	query := fmt.Sprintf("SELECT id, password FROM users WHERE %s = ?", usernameField)
	row := db.QueryRow(query, user.Username)

	var userID int
	var storedPassword string

	err = row.Scan(&userID, &storedPassword)
	if err != nil {
		data["message"] = "Invalid username or password"
		utils.JsonResponse(w, data, 401)
		return
	}

	// Compare stored password hash with provided password
	if !services.CheckPasswordHash(user.Password, storedPassword) {
		data["message"] = "Invalid username or password"
		utils.JsonResponse(w, data, 401)
		return
	}

	// Generate JWT token
	token, errToken := services.GenerateToken(uint(userID))
	if errToken != nil {
		data["message"] = "Problem generating token."
		utils.JsonResponse(w, data, 500)
		return
	}

	// Return success response
	data["message"] = "Login successful."
	data["token"] = token

	utils.JsonResponse(w, data, 200)
}

// ForgotPasswordAuth => Send email or sms to user for forgot password
// @Summary Forgot password
// @Description Send email or sms to user for forgot password
// @Produce json
// @Param username body string true "The email or phone of the user"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/forgot-password [post]
func (db *DBHandler) ForgotPasswordAuth(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username string `json:"username"`
	}

	data := make(map[string]any)

	// Parse body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, 400)
		return
	}

	// Detect username
	usernameField := detectEmailOrPhone(user.Username)

	// Check user is existing
	queryExist := fmt.Sprintf("SELECT count(*) FROM users WHERE %s = ?", usernameField)
	rows, err := db.Query(queryExist, user.Username)
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

	// Check user exists
	if count == 0 {
		data["message"] = "The user is not exists."

		utils.JsonResponse(w, data, 302)
		return
	}

	// Send email or sms for user
	if usernameField == "email" {
		// TODO:

		data["message"] = "Email sent successfully."
	} else {
		// send sms
		// TODO:

		data["message"] = "Sms sent successfully."
	}

	utils.JsonResponse(w, data, 200)
}

// ForgotPasswordVerifyAuth => Verify the user with otp
// @Summary Forgot Password Verify
// @Description Verify the user with otp
// @Produce json
// @Param username body string true "The email or phone of the user"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/forgot-password-verify [post]
func (db *DBHandler) ForgotPasswordVerifyAuth(w http.ResponseWriter, r *http.Request) {
	// TODO:
}

// detectEmailOrPhone => Detects whether the username is an email or phone
func detectEmailOrPhone(username string) string {
	_, err := mail.ParseAddress(username)
	if err == nil {
		return "email"
	}

	// Check if it's a valid phone number (assuming 10+ digits)
	phoneRegex := regexp.MustCompile(`^\d{10,}$`)
	if phoneRegex.MatchString(username) {
		return "phone"
	}

	// Default to phone if unsure
	return "phone"
}
