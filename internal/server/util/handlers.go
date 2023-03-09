package util

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	Error string `json:"error"`
}

func ReadBody(c *gin.Context) (body string, err error) {
	content, err := io.ReadAll(c.Request.Body)
	if err == nil {
		body = string(content)
	}
	return
}

func ValidateContentTypeJSON(c *gin.Context, contentType string) bool {
	headers, ok := c.Request.Header["Content-Type"]
	if !ok {
		response := WriteJSONError(c, "Content-Type must be '"+contentType+"'")
		c.String(http.StatusBadRequest, response)
		return false
	}
	for _, header := range headers {
		if strings.Contains(header, contentType) {
			return true
		}
	}
	response := WriteJSONError(c, "Content-Type must be '"+contentType+"'")
	c.String(http.StatusBadRequest, response)
	return false
}

func WriteJSONError(c *gin.Context, error string) string {
	responseError := ResponseError{Error: error}
	response, _ := json.Marshal(responseError)
	return string(response)
}

func HeaderContains(header http.Header, headerName, headerValue string) bool {
	headers, ok := header[headerName]
	if !ok {
		return false
	}

	for _, h := range headers {
		if strings.Contains(h, headerValue) {
			return true
		}
	}

	return false
}
