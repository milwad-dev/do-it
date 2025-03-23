package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/milwad-dev/do-it/internal/models"
	"github.com/milwad-dev/do-it/internal/utils"
	"net/http"
	"time"
)

// GetLatestTasks => Get the latest tasks and return json response
func (db *DBHandler) GetLatestTasks(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

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
	JOIN labels ON tasks.label_id = labels.id`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		data["message"] = err.Error()
		data["status"] = "error"
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
			data["status"] = "error"
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
		data["status"] = "error"
		utils.JsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	// Add tasks data to response
	data["data"] = tasks
	data["status"] = "success"

	// Respond with the tasks data as JSON
	utils.JsonResponse(w, data, http.StatusOK)
}

// StoreTask => Store new task and return json response
func (db *DBHandler) StoreTask(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	r.ParseForm()

	// TODO: ADD validation

	// Read request body
	var task models.Task

	userId := 9 // TODO FIX THIS

	// Decode JSON request body into `tasks`
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Exec query (id, created_at and updated_at filled automatically)
	query := "INSERT INTO tasks (title, description, status, label_id, user_id) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(query, &task.Title, &task.Description, &task.Status, &task.LabelId, userId)
	if err != nil {
		panic(err)
	}

	data := make(map[string]string)
	data["message"] = "The task store successfully."

	utils.JsonResponse(w, data, 200)
}

// MarkTaskAsCompleted => Mark task as completed
func (db *DBHandler) MarkTaskAsCompleted(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "task")

	sql := "UPDATE tasks SET completed_at = ? WHERE id = ?"
	res, err := db.Exec(sql, time.Now(), taskId)
	if err != nil {
		panic(err)
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		panic(err)
	}

	data := make(map[string]string)
	data["message"] = "The task mark as completed."

	utils.JsonResponse(w, data, 200)
}
