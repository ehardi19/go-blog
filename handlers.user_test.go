// handlers.user_test.go

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestShowLoginPageAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	r.GET("/u/login", ensureNotLoggedIn(), showLoginPage)

	req, _ := http.NewRequest("GET", "/u/login", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}

func TestShowLoginPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	r.GET("/u/login", ensureNotLoggedIn(), showLoginPage)

	req, _ := http.NewRequest("GET", "/u/login", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Login</title>") > 0

		return statusOK && pageOK
	})
}

func TestLoginAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	r.POST("/u/login", ensureNotLoggedIn(), performLogin)

	loginPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/login", strings.NewReader(loginPayload))
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}

func TestLoginUnauthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	r.POST("/u/login", ensureNotLoggedIn(), performLogin)

	loginPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/login", strings.NewReader(loginPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>Successful Login</title>") < 0 {
		t.Fail()
	}
}

func TestLoginUnauthenticatedIncorrectCredentials(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	r.POST("/u/login", ensureNotLoggedIn(), performLogin)

	loginPayload := getRegistrationPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/login", strings.NewReader(loginPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func TestShowRegistrationPageAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	r.GET("/u/register", ensureNotLoggedIn(), showRegistrationPage)

	req, _ := http.NewRequest("GET", "/u/register", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}

func TestShowRegistrationPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	r.GET("/u/register", ensureNotLoggedIn(), showRegistrationPage)

	req, _ := http.NewRequest("GET", "/u/register", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Register</title>") > 0

		return statusOK && pageOK
	})
}

func TestRegisterAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	r.POST("/u/register", ensureNotLoggedIn(), register)

	registrationPayload := getRegistrationPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}

func TestRegisterUnauthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	r.POST("/u/register", ensureNotLoggedIn(), register)

	registrationPayload := getRegistrationPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>Successful registration &amp; Login</title>") < 0 {
		t.Fail()
	}
}

func TestRegisterUnauthenticatedUnavailableUsername(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	r.POST("/u/register", ensureNotLoggedIn(), register)

	registrationPayload := getLoginPOSTPayload()
	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}

func getLoginPOSTPayload() string {
	params := url.Values{}
	params.Add("username", "user1")
	params.Add("password", "pass1")

	return params.Encode()
}

func getRegistrationPOSTPayload() string {
	params := url.Values{}
	params.Add("username", "u1")
	params.Add("password", "p1")

	return params.Encode()
}
