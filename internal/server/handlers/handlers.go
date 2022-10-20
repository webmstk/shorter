package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/storage"
)

func HandlerShorten(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
		if !validateContentType(c, "text/plain") {
			return
		}
		longURL, err := readBody(c)
		if err != nil || longURL == "" {
			c.String(http.StatusBadRequest, "Please specify valid url in body")
			return
		}

		shortURL, err := storage.SaveLongURL(longURL)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to generate short link")
			return
		}

		c.String(http.StatusCreated, config.ServerBaseURL+"/"+shortURL)
	}
}

func HandlerExpand(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
		shortURL := c.Param("shortURL")
		longURL, ok := storage.GetLongURL(shortURL)
		if ok {
			c.Redirect(http.StatusTemporaryRedirect, longURL)
		} else {
			c.String(http.StatusNotFound, "Short url not found")
		}
	}
}

func validateContentType(c *gin.Context, contentType string) bool {
	headers, ok := c.Request.Header["Content-Type"]
	if !ok {
		c.String(http.StatusBadRequest, "Content-Type must be '"+contentType+"'")
		return false
	}
	for _, header := range headers {
		if strings.Contains(header, contentType) {
			return true
		}
	}
	c.String(http.StatusBadRequest, "Content-Type must be '"+contentType+"'")
	return false
}

func readBody(c *gin.Context) (body string, err error) {
	defer func() { _ = c.Request.Body.Close() }()
	content, err := io.ReadAll(c.Request.Body)
	if err == nil {
		body = string(content)
	}
	return
}
