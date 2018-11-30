package main

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func isUserValid(username, password string) bool {
	u := user{Username: username, Password: password}
	result := db.QueryRow("SELECT password FROM userTable WHERE username=$1", u.Password)

	tmp := &u
	err := result.Scan(&tmp.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(tmp.Password), []byte(u.Password)); err != nil {
		return false
	}

	return true
}


func registerNewUser(username, password string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	u := user{Username: username, Password: password}

	return &u, nil
}

func isUsernameAvailable(username string) bool {
	stmt := "SELECT username FROM userTable WHERE username = ?"
	err := db.QueryRow(stmt, username).Scan(&username)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return false
	}
	return true
}

func showLoginPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

func performLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if isUserValid(username, password) {
		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)

		render(c, gin.H{
			"title": "Successful Login"}, "login-successful.html")

	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
	}
}

func generateSessionToken() string {
	return strconv.FormatInt(rand.Int63(), 16)
}

func logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func showRegistrationPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Register"}, "register.html")
}


func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	u := user{Username: username, Password: password}

	if _, err := registerNewUser(username, password); err == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)

		_, err = db.Query("INSERT INTO userTable VALUES ($1, $2)", u.Username, string(hashedPassword))
		if err != nil {
			panic(err)
		}


		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)

		render(c, gin.H{
			"title": "Successful registration & Login"}, "login-successful.html")

	} else {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"ErrorTitle":   "Registration Failed",
			"ErrorMessage": err.Error()})

	}
}



