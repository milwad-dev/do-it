package handlers

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/milwad-dev/do-it/internal/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetLatestLabels_OK(t *testing.T) {
	// Mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Initialize logger in non-production mode for testing (prints JSON to stdout + file)
	logger.InitLogger(false)

	// Expectations
	mock.ExpectQuery(`SELECT l\.id, l\.title, l\.color, l\.created_at, l\.updated_at, l\.user_id, 
\s*u\.id, u\.name, COALESCE\(u\.email, ''\), COALESCE\(u\.phone, ''\), u\.created_at 
\s*FROM labels l
\s*JOIN users u ON l\.user_id = u\.id
\s*WHERE l\.user_id = \?
\s*ORDER BY l\.created_at DESC`).
		WithArgs(float64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "title", "color", "created_at", "updated_at", "user_id",
			"user_id", "name", "email", "phone", "user_created_at",
		}).AddRow(
			1, "Label Title", "#FF0000", "2025-01-01 10:00:00", "2025-01-02 10:00:00", 42,
			42, "User Name", "user@example.com", "1234567890", "2024-12-31 09:00:00",
		))

	// Handler
	h := &DBHandler{DB: db}

	// Request setup
	req := httptest.NewRequest(http.MethodGet, "/labels", nil)

	// Add Chi route context with id param
	req = callContext(req)

	// Recorder & handler call
	rr := httptest.NewRecorder()
	h.GetLatestLabels(rr, req)

	// Assert
	require.Equal(t, http.StatusOK, rr.Code)
	require.JSONEq(t, `{"data":{
{
"color":"#FF0000",
"created_at":"2025-01-01 10:00:00",
"id":1,
"title":"Label Title",
"updated_at":"2025-01-02 10:00:00", 
"user":{
"created_at":"2024-12-31 09:00:00",
"email":"user@example.com",
"emailVerified_at":"0001-01-01T00:00:00Z",
"id":42,
"name":"User Name",
"phone":"1234567890",
"phoneVerified_at":"0001-01-01T00:00:00Z",
"updated_at":""}
}},
"status": "Success"
}`, rr.Body.String())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteLabel_OK(t *testing.T) {
	// Mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Expectations
	mock.ExpectQuery(`SELECT count\(\*\) FROM labels WHERE id = \? AND user_id = \?`).
		WithArgs("42", float64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectExec(`DELETE FROM labels WHERE id = \? AND user_id = \?`).
		WithArgs("42", float64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Handler
	h := &DBHandler{DB: db}

	// Request setup
	req := httptest.NewRequest(http.MethodDelete, "/labels/42", nil)

	// Add Chi route context with id param
	req = callContext(req)

	// Recorder & handler call
	rr := httptest.NewRecorder()
	h.DeleteLabel(rr, req)

	// Assert
	require.Equal(t, http.StatusOK, rr.Code)
	require.JSONEq(t, `{"data": {
"message":"The label deleted successfully."
},
"status": "Success"
}`, rr.Body.String())
	require.NoError(t, mock.ExpectationsWereMet())
}

func callContext(req *http.Request) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "42")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Add JWT claims to context
	req = req.WithContext(context.WithValue(req.Context(), "userID", jwt.MapClaims{
		"user_id": float64(1),
	}))
	return req
}
