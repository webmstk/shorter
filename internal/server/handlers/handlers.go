package handlers

import (
	"fmt"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/storage"
	"io"
	"net/http"
	"strings"
)

func HandlerLinks(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathChunks := splitURLPath(r.URL.Path)
		pathLen := len(pathChunks)

		if pathLen == 0 && r.Method == http.MethodPost {
			// POST /
			HandlerShorten(storage, w, r)
		} else if pathLen == 1 && r.Method == http.MethodGet {
			// GET /123, GET /123/
			HandlerExpand(storage, w, r)
		} else {
			HandlerNotFound(w, r)
		}
	}
}

func HandlerNotFound(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

func HandlerShorten(storage storage.Storage, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if !validateContentType(w, r, "text/plain; charset=utf-8") {
		return
	}
	longURL, err := readBody(r)
	if err != nil || longURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Please specify valid url in body"))
		return
	}

	shortURL, err := storage.SaveLongURL(longURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to generate short link"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(config.ServerBaseURL + "/" + shortURL))
}

func HandlerExpand(storage storage.Storage, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	shortURL := splitURLPath(r.URL.Path)[0]
	longURL, ok := storage.GetLongURL(shortURL)
	if ok {
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Short url not found"))
	}
}

func splitURLPath(path string) (fragments []string) {
	// "/abc/def/" превращается в ["", "abc", "def", ""]
	// поэтому дополнительно чистим от пустых строк
	for _, fragment := range strings.Split(path, "/") {
		if fragment != "" {
			fragments = append(fragments, fragment)
		}
	}
	return
}

func validateContentType(w http.ResponseWriter, r *http.Request, contentType string) (ok bool) {
	ct := r.Header.Get("Content-Type")
	if ct != contentType {
		w.WriteHeader(http.StatusBadRequest)
		body := fmt.Sprintf("Content-Type must be '%s'", contentType)
		_, _ = w.Write([]byte(body))
		return false
	}
	return true
}

func readBody(r *http.Request) (body string, err error) {
	defer func() { _ = r.Body.Close() }()
	content, err := io.ReadAll(r.Body)
	if err == nil {
		body = string(content)
	}
	return
}
