package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/api"
	"example.com/storage"
	"github.com/stretchr/testify/suite"
)

const incorrectLoginPayload = `{
	"email": "incorrect@email.com",
	"password": "wrong"
}`

const loginPayload = `{
	"email": "testemail@email.com",
	"password": "boom"
}`

const accountCratePayload = `{
	"email": "newaccount@email.com",
	"firstName": "Attila",
	"lastName": "Olbrich"
  }`

type Suite struct {
	suite.Suite
}

func TestRun(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (t *Suite) TestIfWorks() {
	t.True(true)
}

func initApp(withSeed bool) *api.APIServer {
	var app api.APIServer
	store, _ := storage.NewDatabaseStore(true)
	if withSeed {
		seedDatabase(store, true)
	}
	store.Init()

	app = api.APIServer{Store: store}
	app.Run()

	return &app
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app := initApp(true)
	app.Router.ServeHTTP(rr, req)

	return rr
}

func executeAuthenticatedRequest(req *http.Request) *httptest.ResponseRecorder {
	token := getToken()
	req.Header.Set("x-jwt-token", token)
	req.Header.Set("Content-type", "application/json")

	rr := httptest.NewRecorder()
	app := initApp(false)
	app.Router.ServeHTTP(rr, req)

	return rr
}

func getToken() string {
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(loginPayload)))
	response := executeRequest(req)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	token, _ := m["token"]

	return token.(string)
}

func (t *Suite) TestIfReturnsPage() {
	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)

	t.Equal(http.StatusOK, response.Code)
	t.Contains(string(response.Body.Bytes()), "<html")
}

func (t *Suite) TestIfLoginReturnsError() {
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(incorrectLoginPayload)))
	response := executeRequest(req)

	t.Equal(http.StatusBadRequest, response.Code)
}

func (t *Suite) TestIfLoginReturnsSuccessfull() {
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(loginPayload)))
	response := executeRequest(req)

	t.Equal(http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	token, ok := m["token"]
	t.True(ok)
	t.Len(token, 125)
}

func (t *Suite) TestIfAccountReturnsForbidden() {
	req, _ := http.NewRequest("GET", "/account", nil)
	response := executeRequest(req)

	t.Equal(http.StatusForbidden, response.Code)
}

func (t *Suite) TestIfAccountReturnsOkIfHasAccessCode() {
	token := getToken()
	req, _ := http.NewRequest("GET", "/account", nil)
	req.Header.Set("x-jwt-token", token)
	req.Header.Set("Content-type", "application/json")

	response := executeRequest(req)

	t.Equal(http.StatusOK, response.Code)
}

func (t *Suite) TestIfAccountCanBeCreated() {
	req, _ := http.NewRequest("POST", "/account", bytes.NewBuffer([]byte(accountCratePayload)))

	response := executeAuthenticatedRequest(req)

	t.Equal(http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	id, ok := m["id"]
	t.True(ok)
	t.Greater(int(id.(float64)), 0)

	email, ok := m["email"]
	t.True(ok)
	t.Equal("newaccount@email.com", email)
}

// @TODO add further tests
