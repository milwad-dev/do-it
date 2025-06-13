package handlers

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLatestUsers_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("mocking error: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at"}).
		AddRow(1, "Alice", "alice@example.com", "123456", "2024-01-01T00:00:00Z").
		AddRow(2, "Bob", "bob@example.com", "654321", "2024-01-02T00:00:00Z")

	mock.ExpectQuery("SELECT id, name, COALESCE").WillReturnRows(rows)

	handler := &DBHandler{DB: db}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler.GetLatestUsers(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	if _, ok := response["data"]; !ok {
		t.Error("expected data in response")
	}
}
