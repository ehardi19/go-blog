// handlers.article_test.go

package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestShowIndexPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	r.GET("/", showIndexPage)

	req, _ := http.NewRequest("GET", "/", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Home Page</title>") > 0

		return statusOK && pageOK
	})
}

func TestShowIndexPageAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	r.GET("/", showIndexPage)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>Home Page</title>") < 0 {
		t.Fail()
	}

}

func TestArticleUnauthenticated(t *testing.T) {
	r := getRouter(true)

	r.GET("/article/view/:article_id", getArticle)

	req, _ := http.NewRequest("GET", "/article/view/1", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "<title>Article 1</title>") > 0

		return statusOK && pageOK
	})
}

func TestArticleAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	r.GET("/article/view/:article_id", getArticle)

	req, _ := http.NewRequest("GET", "/article/view/1", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>Article 1</title>") < 0 {
		t.Fail()
	}

}

func TestArticleListJSON(t *testing.T) {
	r := getRouter(true)

	r.GET("/", showIndexPage)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "application/json")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			return false
		}
		var articles []article
		err = json.Unmarshal(p, &articles)

		return err == nil && len(articles) >= 2 && statusOK
	})
}

func TestArticleXML(t *testing.T) {
	r := getRouter(true)

	r.GET("/article/view/:article_id", getArticle)

	req, _ := http.NewRequest("GET", "/article/view/1", nil)
	req.Header.Add("Accept", "application/xml")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			return false
		}
		var a article
		err = xml.Unmarshal(p, &a)

		return err == nil && a.ID == 1 && len(a.Title) >= 0 && statusOK
	})
}

func TestArticleCreationPageAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	r.GET("/article/create", ensureLoggedIn(), showArticleCreationPage)

	req, _ := http.NewRequest("GET", "/article/create", nil)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>Create New Article</title>") < 0 {
		t.Fail()
	}

}

func TestArticleCreationPageUnauthenticated(t *testing.T) {
	r := getRouter(true)

	r.GET("/article/create", ensureLoggedIn(), showArticleCreationPage)

	req, _ := http.NewRequest("GET", "/article/create", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusUnauthorized
	})
}

func TestArticleCreationAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()

	r := getRouter(true)

	http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

	r.POST("/article/create", ensureLoggedIn(), createArticle)

	articlePayload := getArticlePOSTPayload()
	req, _ := http.NewRequest("POST", "/article/create", strings.NewReader(articlePayload))
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(articlePayload)))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	p, err := ioutil.ReadAll(w.Body)
	if err != nil || strings.Index(string(p), "<title>Submission Successful</title>") < 0 {
		t.Fail()
	}

}

func TestArticleCreationUnauthenticated(t *testing.T) {
	r := getRouter(true)

	r.POST("/article/create", ensureLoggedIn(), createArticle)

	articlePayload := getArticlePOSTPayload()
	req, _ := http.NewRequest("POST", "/article/create", strings.NewReader(articlePayload))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(articlePayload)))

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusUnauthorized
	})
}

func getArticlePOSTPayload() string {
	params := url.Values{}
	params.Add("title", "Test Article Title")
	params.Add("content", "Test Article Content")

	return params.Encode()
}
