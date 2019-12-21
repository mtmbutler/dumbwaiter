package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	// "os"
	"testing"
)

type Case struct {
	serverFunc      http.HandlerFunc
	endpoint        string
	reqType         string
	body            []byte
	params          map[string]string
	expectedCode    int
	expectedBodyStr string
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
		t.Errorf("handler returned unexpected body:\n got `%v`\nwant `%v`",
			rr.Body.String(), c.expectedBodyStr)
	}
}

func TestReturnAllUsersNoAuth(t *testing.T) {
	connectDB()
	defer DB.Close()

	c := Case{
		serverFunc:      returnAllUsers,
		endpoint:        "/users",
		reqType:         "GET",
		body:            []byte(``),
		params:          map[string]string{},
		expectedCode:    http.StatusUnauthorized,
		expectedBodyStr: ``,
	}
	c.run(t)
}

func TestReturnAllUsersNonAdmin(t *testing.T) {
	connectDB()
	defer DB.Close()

	// Create a non-admin user so we can verify the non-admin user can't see the admin
	// user when they hit this endpoint. Note that we don't have to create the admin,
	// since that's done automatically when the server starts.
	user := User{
		Email:   "user@gmail.com",
		ApiKey:  "someApiKey",
		IsAdmin: false,
	}
	DB.FirstOrCreate(&user, user)
	c := Case{
		serverFunc:      returnAllUsers,
		endpoint:        "/users",
		reqType:         "GET",
		body:            []byte(``),
		params:          map[string]string{"apiKey": user.ApiKey},
		expectedCode:    http.StatusOK,
		expectedBodyStr: fmt.Sprintf(`[{"id":%v,"email":"user@gmail.com","apiKey":"someApiKey","isAdmin":false}]`+"\n", user.ID),
	}
	c.run(t)
	if user.ID != 0 {
		DB.Delete(&user)
	}
}

func TestReturnAllUsersAdmin(t *testing.T) {
	connectDB()
	defer DB.Close()

	// Get admin credentials
	var admin User
	DB.First(&admin, User{IsAdmin: true})

	// Create a non-admin user so we can verify the admin user can see the non-admin
	// user when they hit this endpoint. Note that we don't have to create the admin,
	// since that's done automatically when the server starts.
	user := User{
		Email:   "user@gmail.com",
		ApiKey:  "someApiKey",
		IsAdmin: false,
	}
	DB.FirstOrCreate(&user, user)
	c := Case{
		serverFunc:      returnAllUsers,
		endpoint:        "/users",
		reqType:         "GET",
		body:            []byte(``),
		params:          map[string]string{"apiKey": admin.ApiKey},
		expectedCode:    http.StatusOK,
		expectedBodyStr: fmt.Sprintf(`[{"id":%v,"email":"%s","apiKey":"%s","isAdmin":true},{"id":%v,"email":"user@gmail.com","apiKey":"someApiKey","isAdmin":false}]`+"\n", admin.ID, admin.Email, admin.ApiKey, user.ID),
	}
	c.run(t)
	if user.ID != 0 {
		DB.Delete(&user)
	}
}
