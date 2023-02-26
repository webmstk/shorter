package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/server/middlewares"
	"github.com/webmstk/shorter/internal/storage"
)

func TestHandlerShorten(t *testing.T) {
	setupTestConfig(&config.Config)

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
				body:        config.Config.BaseURL + "/" + generateShortLink("https://ya.ru"),
			},
		},
		{
			name:        "same valid link",
			contentType: "text/plain",
			body:        "https://ya.ru",
			want: want{
				contentType: "text/plain",
				statusCode:  201,
				body:        config.Config.BaseURL + "/" + generateShortLink("https://ya.ru"),
			},
		},
		{
			name:        "second valid link",
			contentType: "text/plain",
			body:        "https://yandex.ru",
			want: want{
				contentType: "text/plain",
				statusCode:  201,
				body:        config.Config.BaseURL + "/" + generateShortLink("https://yandex.ru"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupServer(nil)
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
	setupTestConfig(&config.Config)

	linksStorage := storage.NewStorage()
	shortURL, _ := linksStorage.SaveLongURL("https://yandex.ru")

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

func TestCompression(t *testing.T) {
	setupTestConfig(&config.Config)

	longURL := "https://ya.ru"
	linksStorage := storage.NewStorage()
	shortURL, _ := linksStorage.SaveLongURL(longURL)
	data, _ := middlewares.Compress([]byte(longURL))

	t.Run("gzip compression", func(t *testing.T) {
		r := setupServer(linksStorage)

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(data)))
		request.Header.Set("Content-Type", "application/x-gzip")
		request.Header.Set("Content-Encoding", "gzip")
		request.Header.Set("Accept-Encoding", "gzip")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		assert.Equal(t, http.StatusCreated, w.Code)
		// decoded, _ := middlewares.Decompress(w.Body.Bytes())
		// assert.Equal(t, config.Config.BaseURL+"/"+shortURL, string(decoded))
		// assert.Equal(t, "application/x-gzip", w.Header().Get("Content-Type"))
		assert.Equal(t, config.Config.BaseURL+"/"+shortURL, w.Body.String())
	})
}
