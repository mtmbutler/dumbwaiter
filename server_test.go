package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	// "os"
	"testing"
)

type Case struct {
	endpoint        string
	reqType         string
	body            []byte
	params          map[string]string
	expectedCode    int
	expectedBodyStr string
	serverFunc      http.HandlerFunc
}

func (c Case) run(t *testing.T) {
	// Set up the request
	req, err := http.NewRequest(c.reqType, c.endpoint, bytes.NewBuffer(c.body))
	if err != nil {
		t.Fatal(err)
	}

	// Add GET parameters
	q := req.URL.Query()
	for k, v := range c.params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// Set up the handler
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.serverFunc)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != c.expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, c.expectedCode)
	}

	// Check the response body is what we expect.
	if rr.Body.String() != c.expectedBodyStr {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), c.expectedBodyStr)
	}
}

func TestReturnAllUsers(t *testing.T) {
	connectDB()
	defer DB.Close()

	c := Case{
		"/users",
		"GET",
		[]byte(``),
		map[string]string{},
		http.StatusUnauthorized,
		``,
		returnAllUsers,
	}
	c.run(t)
}
