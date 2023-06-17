package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/storage"
)

func TestHandlerAPIShorten(t *testing.T) {
	linksStorage := storage.NewStorage()

	type want struct {
		contentType string
		statusCode  int
		body        string
	}

	tests := []struct {
		name        string
		request     string
		body        string
		contentType string
		want        want
	}{
		{
			name:        "correct request",
			request:     "/api/shorten",
			body:        `{"url":"https://ya.ru"}`,
			contentType: "application/json",
			want: want{
				contentType: "application/json",
				statusCode:  201,
				body:        `{"result":"` + absoluteLink(generateShortLink("https://ya.ru")) + `"}`,
			},
		},
		{
			name:        "request with wrong contentType",
			request:     "/api/shorten",
			body:        `{"url":"https://ya.ru"}`,
			contentType: "text/plain",
			want: want{
				contentType: "application/json",
				statusCode:  400,
				body:        `{"error":"Content-Type must be 'application/json'"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if config.Config.DatabaseDSN != "" {
				db := storage.NewStorageDB()
				db.DeleteLink("https://ya.ru")
			}
			r := setupServer(linksStorage)
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}

func TestHandlerAPIUserLinks(t *testing.T) {
	r := setupServer(nil)

	longURL := "http://ya.ru"
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(longURL))
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

	type want struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name      string
		userID    *http.Cookie
		userToken *http.Cookie
		want      want
	}{
		{
			name:      "with valid auth",
			userID:    userID,
			userToken: userToken,
			want: want{
				statusCode: 200,
				body:       `[{"original_url":"` + longURL + `","short_url":"` + absoluteLink(generateShortLink(longURL)) + `"}]`,
			},
		},
		// {
		// 	name:      "with invalid auth",
		// 	userID:    &http.Cookie{},
		// 	userToken: &http.Cookie{},
		// 	want: want{
		// 		statusCode: 204,
		// 		body:       "",
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			request.AddCookie(tt.userID)
			request.AddCookie(tt.userToken)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}

func TestHandlerAPIBatch(t *testing.T) {
	t.Run("test ping", func(t *testing.T) {
		r := setupServer(nil)
		body := `[
  { "correlation_id": "111", "original_url": "http://rebro.ru" },
  { "correlation_id": "222", "original_url": "http://reshka.ru" }
]`
		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		expectedBody := []storage.BatchOutput{
			{CorrelationID: "111", ShortURL: absoluteLink("886043032")},
			{CorrelationID: "222", ShortURL: absoluteLink("2657218682")},
		}

		var actualBody []storage.BatchOutput
		json.Unmarshal(w.Body.Bytes(), &actualBody)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, expectedBody, actualBody)
	})
}
