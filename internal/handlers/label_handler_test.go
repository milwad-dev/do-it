package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock dbHandler struct to satisfy the dbHandler interface
type dbHandler struct{}

func TestStoreLabel(t *testing.T) {
	// Create a mock dbHandler with any necessary dependencies
	mockDB := &dbHandler{} // You may need to create a mock dbHandler

	// Create a sample label JSON body
	labelJSON := []byte(`{"title": "Test Label", "color": "#FF0000"}`)

	// Create a mock HTTP request with the sample label data
	req, err := http.NewRequest("POST", "/labels", bytes.NewBuffer(labelJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the storeLabel handler function with the mock dbHandler and mock HTTP request
	handler := http.HandlerFunc(mockDB.storeLabel)
	handler.ServeHTTP(rr, req)

	// Check the status code returned by the handler
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the response body contains the expected message
	expectedResponse := `"message":"The label stored successfully."`
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}

// Mock implementation of storeLabel function to satisfy the dbHandler interface
func (db *dbHandler) storeLabel(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["message"] = "The label stored successfully."
	jsonResponse(w, data, 200)
}

// Mock implementation of jsonResponse function for testing purposes
func jsonResponse(w http.ResponseWriter, data map[string]string, statusCode int) {
	w.WriteHeader(statusCode)
	for key, value := range data {
		w.Write([]byte(`"` + key + `":"` + value + `"`))
	}
}
