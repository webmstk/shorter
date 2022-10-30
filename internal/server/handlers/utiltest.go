package handlers

import "github.com/webmstk/shorter/internal/storage"

func generateShortLink(s string) string {
	s, err := storage.GenerateShortLink(s)
	if err != nil {
		panic("Failed to generate short link")
	}
	return s
}
