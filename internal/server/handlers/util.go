package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

type ResponseError struct {
	Error string `json:"error"`
}

func readBody(c *gin.Context) (body string, err error) {
	defer func() { _ = c.Request.Body.Close() }()
	content, err := io.ReadAll(c.Request.Body)
	if err == nil {
		body = string(content)
	}
	return
}

func validateContentTypeJSON(c *gin.Context, contentType string) bool {
	headers, ok := c.Request.Header["Content-Type"]
	if !ok {
		response := writeJSONError(c, "Content-Type must be '"+contentType+"'")
		c.String(http.StatusBadRequest, response)
		return false
	}
	for _, header := range headers {
		if strings.Contains(header, contentType) {
			return true
		}
	}
	response := writeJSONError(c, "Content-Type must be '"+contentType+"'")
	c.String(http.StatusBadRequest, response)
	return false
}

func writeJSONError(c *gin.Context, error string) string {
	responseError := ResponseError{Error: error}
	response, _ := json.Marshal(responseError)
	return string(response)
}
