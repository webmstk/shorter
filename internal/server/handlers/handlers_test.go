package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerShorten(t *testing.T) {
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
			name:        "wrong content type",
			contentType: "wrong_type",
			body:        "https://ya.ru",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "Content-Type must be 'text/plain; charset=utf-8'",
			},
		},
		{
			name:        "empty body",
			contentType: "text/plain; charset=utf-8",
			body:        "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				body:        "Please specify valid url in body",
			},
		},
		{
			name:        "valid link",
			contentType: "text/plain; charset=utf-8",
			body:        "https://ya.ru",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				body:        config.ServerBaseURL + "/" + generateShortLink("https://ya.ru"),
			},
		},
		{
			name:        "same valid link",
			contentType: "text/plain; charset=utf-8",
			body:        "https://ya.ru",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				body:        config.ServerBaseURL + "/" + generateShortLink("https://ya.ru"),
			},
		},
		{
			name:        "second valid link",
			contentType: "text/plain; charset=utf-8",
			body:        "https://yandex.ru",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				body:        config.ServerBaseURL + "/" + generateShortLink("https://yandex.ru"),
			},
		},
	}

	linksStorage := storage.NewStorage()
	gin.SetMode(gin.ReleaseMode)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := SetupRouter(linksStorage)
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)

			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}

func TestHandlerExpand(t *testing.T) {
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
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
				body:        "Short url not found",
			},
		},
		{
			name:    "existing link",
			request: "/" + shortURL,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  307,
				body:        "",
			},
		},
	}

	gin.SetMode(gin.ReleaseMode)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := SetupRouter(linksStorage)
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)

			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}

func generateShortLink(s string) string {
	s, err := storage.GenerateShortLink(s)
	if err != nil {
		panic("Failed to generate short link")
	}
	return s
}
