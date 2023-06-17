package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/server/middlewares"
	"github.com/webmstk/shorter/internal/storage"
)

func TestHandlerShorten(t *testing.T) {
	linksStorage := storage.NewStorage()

	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	tests := []struct {
		name        string
		body        string
		contentType string
		want        want
	}{
		{
			name:        "empty body",
			contentType: "text/plain",
			body:        "",
			want: want{
				contentType: "text/plain",
				statusCode:  400,
				body:        "Please specify valid url in body",
			},
		},
		{
			name:        "valid link",
			contentType: "text/plain",
			body:        "https://ya.ru",
			want: want{
				contentType: "text/plain",
				statusCode:  201,
				body:        absoluteLink(generateShortLink("https://ya.ru")),
			},
		},
		{
			name:        "same valid link",
			contentType: "text/plain",
			body:        "https://ya.ru",
			want: want{
				contentType: "text/plain",
				statusCode:  409,
				body:        absoluteLink(generateShortLink("https://ya.ru")),
			},
		},
		{
			name:        "second valid link",
			contentType: "text/plain",
			body:        "https://yandez.ru",
			want: want{
				contentType: "text/plain",
				statusCode:  201,
				body:        absoluteLink(generateShortLink("https://yandez.ru")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if config.Config.DatabaseDSN != "" && tt.name != "same valid link" {
				db := storage.NewStorageDB()
				db.DeleteLink(tt.body)
			}
			r := setupServer(linksStorage)
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}

func TestHandlerExpand(t *testing.T) {
	linksStorage := storage.NewStorage()
	shortURL, _ := linksStorage.SaveLongURL("https://yandex.ru", "")

	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "not existing link",
			request: "/short_link404",
			want: want{
				contentType: "text/plain",
				statusCode:  404,
				body:        "Short url not found",
			},
		},
		{
			name:    "existing link",
			request: "/" + shortURL,
			want: want{
				contentType: "text/plain",
				statusCode:  307,
				body:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupServer(linksStorage)
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}

func TestHandlerPing(t *testing.T) {
	t.Run("test ping", func(t *testing.T) {
		r := setupServer(nil)
		request := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandlerShortenCookie(t *testing.T) {
	linksStorage := storage.NewStorage()
	r := setupServer(linksStorage)

	t.Run("with no cookies", func(t *testing.T) {
		if config.Config.DatabaseDSN != "" {
			db := storage.NewStorageDB()
			db.DeleteLink("http://yac.ru")
		}
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://yac.ru"))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)
		result := w.Result()
		defer result.Body.Close()

		var userID *http.Cookie
		var userToken *http.Cookie

		for _, cookie := range result.Cookies() {
			if cookie.Name == "user_id" {
				userID = cookie
			}
			if cookie.Name == "user_token" {
				userToken = cookie
			}
		}
		assert.NotNil(t, userID)
		assert.NotNil(t, userToken)
	})

	t.Run("with valid cookie", func(t *testing.T) {
		if config.Config.DatabaseDSN != "" {
			db := storage.NewStorageDB()
			db.DeleteLink("http://yab.ru")
		}
		user := linksStorage.CreateUser()
		cookieID := &http.Cookie{
			Name:  "user_id",
			Value: user,
		}
		signed := signCookie(user)
		cookieToken := &http.Cookie{
			Name:  "user_token",
			Value: signed,
		}
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://yab.ru"))
		request.AddCookie(cookieID)
		request.AddCookie(cookieToken)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)
		result := w.Result()
		defer result.Body.Close()

		var userID *http.Cookie
		var userToken *http.Cookie

		for _, cookie := range result.Cookies() {
			if cookie.Name == "user_id" {
				userID = cookie
			}
			if cookie.Name == "user_token" {
				userToken = cookie
			}
		}
		require.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, userID.Value, cookieID.Value)
		assert.Equal(t, userToken.Value, cookieToken.Value)
	})

	t.Run("with invalid cookie", func(t *testing.T) {
		if config.Config.DatabaseDSN != "" {
			db := storage.NewStorageDB()
			db.DeleteLink("http://yaa.ru")
		}

		cookieID := &http.Cookie{
			Name:  "user_id",
			Value: "123",
		}
		cookieToken := &http.Cookie{
			Name:  "user_token",
			Value: "wrong_token",
		}
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://yaa.ru"))
		request.AddCookie(cookieID)
		request.AddCookie(cookieToken)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)
		result := w.Result()
		defer result.Body.Close()

		var userID *http.Cookie
		var userToken *http.Cookie
		for _, cookie := range result.Cookies() {
			if cookie.Name == "user_id" {
				userID = cookie
			}
			if cookie.Name == "user_token" {
				userToken = cookie
			}
		}
		assert.NotEqual(t, userID.Value, cookieID.Value)
		assert.NotEqual(t, userToken.Value, cookieToken.Value)
	})
}

func TestCompression(t *testing.T) {
	longURL := "https://yad.ru"
	linksStorage := storage.NewStorage()
	shortURL, _ := linksStorage.SaveLongURL(longURL, "")
	data, _ := middlewares.Compress([]byte(longURL))

	t.Run("gzip compression", func(t *testing.T) {
		r := setupServer(linksStorage)

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(data)))
		request.Header.Set("Content-Type", "application/x-gzip")
		request.Header.Set("Content-Encoding", "gzip")
		request.Header.Set("Accept-Encoding", "gzip")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		assert.Equal(t, http.StatusConflict, w.Code)
		decoded, _ := middlewares.Decompress(w.Body.Bytes())
		assert.Equal(t, absoluteLink(shortURL), string(decoded))
		assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	})
}
