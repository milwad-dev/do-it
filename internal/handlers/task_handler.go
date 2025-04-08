package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/milwad-dev/do-it/internal/models"
	"github.com/milwad-dev/do-it/internal/repositories"
	"github.com/milwad-dev/do-it/internal/utils"
	"net/http"
	"time"
)

// GetLatestTasks => Get the latest tasks and return json response
// @Summary Get Tasks
// @Description Get the latest tasks
// @Produce json
// @Success 200 {object} []models.Task
// @Router /tasks [get]
func (db *DBHandler) GetLatestTasks(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	// Get user id from context
	userId := repositories.GetUserIdFromContext(r)

	// SQL query to join tasks and users
	query := `SELECT 
	   	tasks.id AS task_id,
		tasks.title AS task_title,
		tasks.description AS task_description, 
		tasks.status AS task_status, 
		tasks.user_id, 
		tasks.label_id, 
		COALESCE(tasks.completed_at, '') AS task_completed_at,
		tasks.created_at AS task_created_at,
		
		users.id AS user_id, 
		users.name AS user_name,
		COALESCE(users.email, '') AS user_email,
		COALESCE(users.phone, '') AS user_phone,
		users.created_at AS user_created_at,
	
		labels.id AS label_id,
		labels.title AS label_title,
		labels.color AS label_color,
		labels.created_at AS label_created_at
	FROM tasks
	JOIN users ON tasks.user_id = users.id
	JOIN labels ON tasks.label_id = labels.id
	WHERE tasks.user_id = ?`

	// Execute the query
	rows, err := db.Query(query, userId)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Slice to store tasks
	var tasks []models.Task

	// Iterate over the rows
	for rows.Next() {
		var task models.Task
		var user models.User
		var label models.Label

		// Scan the task and user fields from the query result
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.UserId,
			&task.LabelId,
			&task.CompletedAt,
			&task.CreatedAt,

			// Scan user fields
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.CreatedAt,

			// Scan label fields
			&label.ID,
			&label.Title,
			&label.Color,
			&label.CreatedAt,
		)

		// Check for any errors in scanning
		if err != nil {
			data["message"] = "Error scanning row: " + err.Error()

			utils.JsonResponse(w, data, http.StatusInternalServerError)
			return
		}

		// Attach the user to the task
		task.User = user

		// Attach the label to the task
		task.Label = label

		// Append the task to the tasks slice
		tasks = append(tasks, task)
	}

	// Check for any row iteration errors
	if err := rows.Err(); err != nil {
		data["message"] = "Error iterating rows: " + err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	// Add tasks data to response
	data["data"] = tasks

	// Respond with the tasks data as JSON
	utils.JsonResponse(w, data, http.StatusOK)
}

// StoreTask => Store new task and return json response
// @Summary Store Task
// @Description store new task
// @Produce json
// @Param title body string true "The title of the task"
// @Param description body string true "The description of the task"
// @Param status body string true "The status of the task"
// @Param label_id body string true "The label ID of the task"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tasks [post]
func (db *DBHandler) StoreTask(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	r.ParseForm()

	data := make(map[string]string)

	// Read request body
	var task models.Task

	// Get userId from context
	userId := repositories.GetUserIdFromContext(r)

	// Decode JSON request body into `tasks`
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Create a new validator instance
	validate := validator.New()

	// Validate the User struct
	err = validate.Struct(task)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusUnprocessableEntity)
		return
	}

	// Exec query (id, created_at and updated_at filled automatically)
	query := "INSERT INTO tasks (title, description, status, label_id, user_id) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(query, &task.Title, &task.Description, &task.Status, &task.LabelId, userId)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	data["message"] = "The task store successfully."

	utils.JsonResponse(w, data, http.StatusOK)
}

// DeleteTask => Delete task by id
// @Summary Delete task
// @Description Delete task by id
// @Produce json
// @Param id query string true "The ID of the task"
// @Success 200 {object} map[string]string
// @Router /tasks/{id} [delete]
func (db *DBHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get task id from url
	taskId := chi.URLParam(r, "id")

	// Get user id from context
	userId := repositories.GetUserIdFromContext(r)

	// Create data
	data := make(map[string]string)

	// Check task exists
	sql := "SELECT count(*) FROM tasks where id = ? AND user_id = ?"
	row := db.QueryRow(sql, taskId, userId)

	var count int

	err := row.Scan(&count)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	if count == 0 {
		data["message"] = "The task not found."

		utils.JsonResponse(w, data, http.StatusNotFound)
		return
	}

	// Delete task from DB
	sql = "DELETE FROM tasks WHERE id = ? AND user_id = ?"
	_, err = db.Exec(sql, taskId, userId)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	data["message"] = "The task deleted successfully."

	utils.JsonResponse(w, data, http.StatusOK)
}

// MarkTaskAsCompleted => Mark task as completed
// @Summary Mark Task as Completed
// @Description mark task as completed
// @Produce json
// @Param id query string true "The ID of the task"
// @Success 200 {object} map[string]string
// @Router /tasks/{id}/mark-as-completed [patch]
func (db *DBHandler) MarkTaskAsCompleted(w http.ResponseWriter, r *http.Request) {
	// Get task id from url
	taskId := chi.URLParam(r, "id")

	// Get user id from context
	userId := repositories.GetUserIdFromContext(r)

	// Create data
	data := make(map[string]string)

	sql := "UPDATE tasks SET completed_at = ? WHERE id = ? AND user_id = ?"
	res, err := db.Exec(sql, time.Now(), taskId, userId)
	if err != nil {
		data["message"] = err.Error()

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	// If not row affected, we're return error
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		data["message"] = "Error on get affected rows."

		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	data["message"] = "The task mark as completed."

	utils.JsonResponse(w, data, http.StatusOK)
}
