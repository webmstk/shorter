package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	defer logRequest(r, &status)

	pathLen := len(splitURLPath(r.URL.Path))

	if pathLen == 0 && r.Method == http.MethodPost {
		// POST /
		createShortLinkAction(w, r, &status)
	} else if pathLen == 1 && r.Method == http.MethodGet {
		// GET /123, GET /123/
		redirectToLongLinkAction(w, r, &status)
	} else {
		notFoundAction(w, r, &status)
	}
}

func notFoundAction(w http.ResponseWriter, r *http.Request, status *int) {
	*status = http.StatusNotFound
	http.Error(w, "404 Not Found", *status)
}

func createShortLinkAction(w http.ResponseWriter, r *http.Request, status *int) {
	// Validate Content-Type
	ct := r.Header.Get("Content-Type")
	if ct != "text/plain" {
		*status = http.StatusBadRequest
		w.WriteHeader(*status)
		w.Write([]byte("Content-Type must be 'text/plain'"))
		return
	}

	// Validate request Body
	var body string
	b, err := io.ReadAll(r.Body)
	if err == nil {
		body = string(b)
	}
	// можно ужесточить валидацию, но надо знать спецификацию
	if err != nil || body == "" {
		*status = http.StatusBadRequest
		w.WriteHeader(*status)
		w.Write([]byte("Please specify valid url in body"))
		return
	}

	// Suscessful response
	*status = http.StatusCreated
	w.WriteHeader(*status)
	w.Write([]byte("http://short.url"))
}

func redirectToLongLinkAction(w http.ResponseWriter, r *http.Request, status *int) {
	*status = http.StatusTemporaryRedirect
	w.Header().Set("Location", "http://long.url")
	w.WriteHeader(*status)
}

func logRequest(r *http.Request, status *int) {
	log.Println(r.Method, r.RequestURI, *status)
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

func main() {
	HOST := "localhost"
	PORT := "8080"

	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         HOST + ":" + PORT,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/", http.HandlerFunc(rootHandler))

	log.Println("Starting web-server at port", PORT)
	err := s.ListenAndServe()
	log.Fatal(err)
}
