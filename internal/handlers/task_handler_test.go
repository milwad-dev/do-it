package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/milwad-dev/do-it/internal/models"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetLatestTasks_Success(t *testing.T) {
	// Setup DB mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	handler := &DBHandler{DB: db}

	// Mock rows
	rows := sqlmock.NewRows([]string{
		"task_id", "task_title", "task_description", "task_status", "user_id", "label_id", "task_completed_at", "task_created_at",
		"user_id", "user_name", "user_email", "user_phone", "user_created_at",
		"label_id", "label_title", "label_color", "label_created_at",
	}).AddRow(
		1, "Test Task", "A task", "open", 10, 100, "", "2023-01-01",
		10, "Milwad", "milwad@example.com", "123456", "2022-12-01",
		100, "Work", "#ffcc00", "2022-11-01",
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
	SELECT 
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
	WHERE tasks.user_id = ?
`)).WithArgs(float64(1)).WillReturnRows(rows)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Content-Type", "application/json")

	// Mock user id from context
	req = callContext(req)

	// Record response
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetLatestTasks(rr, req)

	// Check status
	if rr.Code != http.StatusOK {
		t.Errorf("Error: %s", rr.Body)
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestStoreTask_Success(t *testing.T) {
	// Setup sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	handler := &DBHandler{DB: db}

	// Prepare input task as JSON
	task := models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
		LabelId:     1,
	}
	taskJSON, _ := json.Marshal(task)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	// Mock user id from context
	req = callContext(req)

	// Record response
	w := httptest.NewRecorder()

	// Mock the label exists query returning count = 1
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM labels WHERE id = \\? AND user_id = \\?").
		WithArgs(task.LabelId, float64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock the INSERT query
	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(task.Title, task.Description, task.Status, task.LabelId, float64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function
	handler.StoreTask(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK but got %d", resp.StatusCode)
	}
}
