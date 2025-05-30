package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/milwad-dev/do-it/internal/repositories"
	"github.com/milwad-dev/do-it/internal/services"
	"github.com/milwad-dev/do-it/internal/utils"
	"log"
	"net/http"
	"net/mail"
	"regexp"
)

var ctx = context.Background()

// RegisterAuth => Register user and create token
// @Summary Register user
// @Description Create new user with token
// @Produce json
// @Param name body string true "The name of the user"
// @Param username body string true "The email or phone of the user"
// @Param password body string true "The password of the user"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /register [post]
func (db *DBHandler) RegisterAuth(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Name     string `json:"name" validate:"required,min=3,max=250"`
		Username string `json:"username" validate:"required,min=3,max=250"`
		Password string `json:"password" validate:"required,min=8,max=250"`
	}

	data := make(map[string]any)

	// Parse body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusBadRequest)
		return
	}

	// Create a new validator instance
	validate := validator.New()

	// Validate the User struct
	err = validate.Struct(user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
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
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
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

		utils.JsonResponse(w, data, http.StatusFound)
		return
	}

	// Hash password
	hashPassword, _ := services.HashPassword(password)

	// Create user
	query := fmt.Sprintf("INSERT INTO users (name, %s, password) VALUES (?, ?, ?)", usernameField)
	result, errInsert := db.Exec(query, name, username, hashPassword)
	if errInsert != nil {
		data["message"] = fmt.Sprintf("Problem on creating user: %s", errInsert.Error())

		utils.JsonResponse(w, data, http.StatusFound)
		return
	}

	// Create token
	userId, _ := result.LastInsertId()
	token, errToken := services.GenerateToken(uint(userId))
	if errToken != nil {
		data["message"] = "Problem on generating token."

		utils.JsonResponse(w, data, http.StatusFound)
		return
	}

	// Return success response
	data["message"] = "Register completed."
	data["token"] = token

	utils.JsonResponse(w, data, http.StatusOK)
}

// LoginAuth => Check user credentials and create jwt token
// @Summary Login user
// @Description Check user credentials and login
// @Produce json
// @Param username body string true "The email or phone of the user"
// @Param password body string true "The password of the user"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /login [post]
func (db *DBHandler) LoginAuth(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Username string `json:"username" validate:"required,min=3,max=250"`
		Password string `json:"password" validate:"required,min=8,max=250"`
	}{}
	data := make(map[string]any)

	// Parse body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data["message"] = err.Error()
		utils.JsonResponse(w, data, http.StatusBadRequest)
		return
	}

	// Create a new validator instance
	validate := validator.New()

	// Validate the User struct
	err = validate.Struct(user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
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
		data["message"] = "Invalid username or password."

		utils.JsonResponse(w, data, http.StatusFound)
		return
	}

	// Compare stored password hash with provided password
	if !services.CheckPasswordHash(user.Password, storedPassword) {
		data["message"] = "Invalid username or password"

		utils.JsonResponse(w, data, http.StatusFound)
		return
	}

	// Generate JWT token
	token, errToken := services.GenerateToken(uint(userID))
	if errToken != nil {
		data["message"] = "Problem generating token."

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	// Return success response
	data["message"] = "Login successful."
	data["token"] = token

	utils.JsonResponse(w, data, http.StatusOK)
}

// ForgotPasswordAuth => Send email or sms to user for forgot password
// @Summary Forgot password
// @Description Send email or sms to user for forgot password
// @Produce json
// @Param username body string true "The email or phone of the user"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /forgot-password [post]
func (db *DBHandler) ForgotPasswordAuth(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username string `json:"username" validate:"required,min=3,max=250"`
		Name     string
		Id       int
	}

	data := make(map[string]any)

	// Parse body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusBadRequest)
		return
	}

	// Create a new validator instance
	validate := validator.New()

	// Validate the User struct
	err = validate.Struct(user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
		return
	}

	// Detect username
	usernameField := detectEmailOrPhone(user.Username)

	// Check user is existing
	queryExist := fmt.Sprintf("SELECT count(*) FROM users WHERE %s = ?", usernameField)
	rows, err := db.Query(queryExist, user.Username)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			data["message"] = err.Error()

			utils.JsonResponse(w, data, http.StatusInternalServerError)
			return
		}
	}

	// Check user exists
	if count == 0 {
		data["message"] = "The user is not exists."

		utils.JsonResponse(w, data, http.StatusFound)
		return
	}

	// Read name of the user from db
	row, err := db.Query(fmt.Sprintf("SELECT id, name FROM users WHERE %s = ?", usernameField), user.Username)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	defer row.Close()

	for row.Next() {
		if err := row.Scan(&user.Id, &user.Name); err != nil {
			data["message"] = err.Error()

			utils.JsonResponse(w, data, http.StatusInternalServerError)
			return
		}
	}

	// Generate random code
	code := utils.NumberBetween(1000, 9999)

	// Set a key-value pair
	err = db.redisClient.Set(ctx, "forgot-password-"+string(rune(user.Id)), code, 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	// Send email or sms for user
	if usernameField == "email" {
		mailData := struct {
			Subject string
			Name    string
			Body    string
		}{
			Subject: "Forgot Password",
			Name:    user.Name,
			Body:    "You forgot password",
		}
		err := services.SendMail(mailData, "forgot-password.html", user.Username)
		if err != nil {
			data["message"] = err.Error()

			utils.JsonResponse(w, data, http.StatusInternalServerError)
			return
		}

		data["message"] = "Email sent successfully."
	} else {
		// TODO:

		data["message"] = "Sms sent successfully."
	}

	utils.JsonResponse(w, data, http.StatusOK)
}

// ForgotPasswordVerifyAuth => Verify the user with otp
// @Summary Forgot Password Verify
// @Description Verify the user with otp
// @Produce json
// @Param username body string true "The email or phone of the user"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /forgot-password-verify [post]
func (db *DBHandler) ForgotPasswordVerifyAuth(w http.ResponseWriter, r *http.Request) {
	var user struct {
		id          int
		password    string `json:"password" validate:"required,min=8,max=250"`
		re_password string `json:"re_password" validate:"required,min=8,max=250"`
		code        string `json:"code" validate:"required"`
	}
	user.id = repositories.GetUserIdFromContext(r).(int)
	data := make(map[string]any)

	code, err := db.redisClient.Get(r.Context(), "forgot-password-"+string(rune(user.id))).Result()
	if err != nil {
		data["message"] = "Try again."

		utils.JsonResponse(w, data, http.StatusFound)
		return
	}

	// Parse body
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusBadRequest)
		return
	}

	// Create a new validator instance
	validate := validator.New()

	// Validate the User struct
	err = validate.Struct(user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
		return
	}

	// Check verify code is valid
	if code != user.code {
		data["message"] = "The code is not valid."

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
		return
	}

	// Check password be same
	if user.password != user.re_password {
		data["message"] = "The password is not valid."

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
		return
	}

	// Update user password
	query := "UPDATE users SET password = ? WHERE id = ?"
	password, _ := services.HashPassword(user.password)
	_, err = db.Exec(query, password, user.id)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	data["message"] = "The password update successfully."
	utils.JsonResponse(w, data, http.StatusOK)
}

// ResetPasswordAuth => Reset user password
// @Summary Reset Password
// @Description Reset user password
// @Produce json
// @Param username body string true "The email or phone of the user"
// @Param new_password body string true "The new password"
// @Param re_new_password body string true "The retry new password"
// @Success 200 {object} map[string]string
// @Failure 302 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reset-password [post]
func (db *DBHandler) ResetPasswordAuth(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username      string `json:"username" validate:"required,min=3,max=250"`
		NewPassword   string `json:"new_password" validate:"required,min=8,max=250"`
		ReNewPassword string `json:"re_new_password" validate:"required,min=8,max=250"`
	}

	data := make(map[string]any)

	// Parse body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusBadRequest)
		return
	}

	// Create a new validator instance
	validate := validator.New()

	// Validate the User struct
	err = validate.Struct(user)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
		return
	}

	// Detect username
	usernameField := detectEmailOrPhone(user.Username)

	// Check user is existing
	queryExist := fmt.Sprintf("SELECT count(*) FROM users WHERE %s = ?", usernameField)
	rows, err := db.Query(queryExist, user.Username)
	if err != nil {
		data["message"] = "Problem on generating token."

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			data["message"] = err.Error()

			utils.JsonResponse(w, data, http.StatusInternalServerError)
			return
		}
	}

	// Check user exists
	if count == 0 {
		data["message"] = "The user is not exists."

		utils.JsonResponse(w, data, http.StatusFound)
		return
	}

	// Check `new_password` is exists with `re_new_password`
	if user.NewPassword != user.ReNewPassword {
		data["message"] = "The new password is not the same with retry new password."

		utils.JsonResponse(w, data, http.StatusFound)
		return
	}

	// Update user password
	password, _ := services.HashPassword(user.NewPassword)
	query := fmt.Sprintf("UPDATE users SET password = ? WHERE %s = ?", usernameField)
	_, err = db.Exec(query, password, user.Username)
	if err != nil {
		data["message"] = "Problem on updating password."

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	data["message"] = "Password updated successfully."
	utils.JsonResponse(w, data, http.StatusOK)
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
