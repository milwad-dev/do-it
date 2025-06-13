package handlers

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

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
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "42")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Add JWT claims to context
	req = req.WithContext(context.WithValue(req.Context(), "userID", jwt.MapClaims{
		"user_id": float64(1),
	}))

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
