package handlers

import (
	"github.com/milwad-dev/do-it/internal/models"
	"github.com/milwad-dev/do-it/internal/utils"
	"net/http"
)

func (db *DBHandler) GetLatestTasks(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	query := `SELECT 
    	tasks.id,
    	tasks.title,
    	tasks.description, 
    	tasks.status, 
    	tasks.user_id, 
    	tasks.label_id, 
    	tasks.completed_at,
    	tasks.created_at,
    	
    	users.id, 
    	users.name,
    	users.email,
    	users.created_at
		FROM tasks
		INNER JOIN users ON tasks.user_id = users.id 
	`
	rows, err := db.Query(query)
	if err != nil {
		data["message"] = err.Error()
		data["status"] = "error"

		utils.JsonResponse(w, data, http.StatusInternalServerError)

		return
	}

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		var user models.User

		rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.UserId,
			&task.LabelId,
			&task.CompletedAt,
			&task.CreatedAt,

			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
		) // TODO: Fix problem for load user

		task.User = user

		tasks = append(tasks, task)
	}

	data["data"] = tasks
	data["status"] = "success"

	utils.JsonResponse(w, data, 200)
}

func (db *DBHandler) StoreTask(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	r.ParseForm()

	// TODO: ADD validation

	// Read request body
	var task *models.Task
	utils.ReadRequestBody(w, r, task)

	userId := 1                       // TODO FIX THIS
	labelId := r.Form.Get("label_id") // TODO: ADD validation is valid id or not

	// Exec query (id, created_at and updated_at filled automatically)
	query := "INSERT INTO tasks (title, description, status, label_id, user_id) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, &task.Title, &task.Description, &task.Status, labelId, userId)
	if err != nil {
		panic(err)
	}

	data := make(map[string]string)
	data["message"] = "The task store successfully."

	utils.JsonResponse(w, data, 200)
}
