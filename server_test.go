package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	// "os"
	"testing"
)

func TestReturnAllUsers(t *testing.T) {
	connectDB()
	defer DB.Close()

	jsonStr := []byte(``)

	// Set up the request
	req, err := http.NewRequest("GET", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	// q := req.URL.Query()
	// q.Add("apiKey", os.Getenv("ADMIN_KEY"))
	// req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(returnAllUsers)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body is what we expect.
	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
