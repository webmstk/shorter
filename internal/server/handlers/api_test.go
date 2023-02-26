package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/storage"
)

func TestHandlerApiExpand(t *testing.T) {
	setupTestConfig(&config.Config)
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
				body:        `{"result":"` + config.Config.BaseURL + "/" + generateShortLink("https://ya.ru") + `"}`,
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
