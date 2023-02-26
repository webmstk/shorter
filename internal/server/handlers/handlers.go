package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/server/util"
	"github.com/webmstk/shorter/internal/storage"
)

func HandlerShorten(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")

		longURL, err := util.ReadBody(c)
		if err != nil || longURL == "" {
			c.String(http.StatusBadRequest, "Please specify valid url in body")
			return
		}

		shortURL, err := storage.SaveLongURL(longURL)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to generate short link")
			return
		}

		c.String(http.StatusCreated, config.Config.BaseURL+"/"+shortURL)
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
