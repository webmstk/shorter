package middlewares

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/server/util"
)

type gzipWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (gw gzipWriter) Write(b []byte) (int, error) {
	return gw.buf.Write(b)
}

func Compressor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Decoding request
		if util.HeaderContains(c.Request.Header, "Content-Encoding", "gzip") {
			body, err := util.ReadBody(c)
			if err != nil {
				abortWithError(c, err)
				return
			}

			newBody, err := Decompress([]byte(body))
			if err != nil {
				abortWithError(c, err)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(newBody))
		}

		// Encoding response
		if util.HeaderContains(c.Request.Header, "Accept-Encoding", "gzip") {
			cw := &gzipWriter{buf: &bytes.Buffer{}, ResponseWriter: c.Writer}
			c.Writer = cw

			c.Next()

			body := cw.buf.String()
			newBody, err := Compress([]byte(body))
			if err != nil {
				abortWithError(c, err)
				return
			}

			c.Header("Content-Encoding", "gzip")
			c.Writer.WriteString(string(newBody))

		} else {
			c.Next()
		}
	}
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return b.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}

func abortWithError(c *gin.Context, err error) {
	msg := fmt.Sprintf("%v", err)
	if util.HeaderContains(c.Request.Header, "Content-Type", "application/json") {
		response := util.WriteJSONError(c, msg)
		c.String(http.StatusBadRequest, response)
	} else {
		c.String(http.StatusBadRequest, fmt.Sprintf("%v", err))
	}
	c.Abort()
}
