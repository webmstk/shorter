package handlers

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/server/util"
	"github.com/webmstk/shorter/internal/storage"
)

func HandlerShorten(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")

		longURL, err := util.ReadBody(c)
		if err != nil || longURL == "" {
			c.String(http.StatusBadRequest, "Please specify valid url in body")
			return
		}

		userID := c.GetString("user_id")
		userToken := c.GetString("user_token")

		status := http.StatusCreated
		shortURL, err := store.SaveLongURL(c, longURL, userID)
		if err != nil {
			var linkExistError *storage.LinkExistError
			if errors.As(err, &linkExistError) {
				status = http.StatusConflict
			} else {
				c.String(http.StatusInternalServerError, "Failed to generate short link")
				return
			}
		}

		host := strings.Split(config.Config.ServerAddress, ":")[0]
		c.SetCookie("user_id", userID, config.Config.CookieTTLSeconds, "/", host, false, true)
		c.SetCookie("user_token", userToken, config.Config.CookieTTLSeconds, "/", host, false, true)
		c.String(status, absoluteLink(shortURL))
	}
}

func HandlerExpand(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
		shortURL := c.Param("shortURL")
		longURL, ok := storage.GetLongURL(c, shortURL)
		if ok {
			c.Redirect(http.StatusTemporaryRedirect, longURL)
		} else {
			c.String(http.StatusNotFound, "Short url not found")
		}
	}
}

func HandlerPing(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch s := store.(type) {
		case *storage.StorageDB:
			err := s.Ping(c)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Status(http.StatusOK)
		default:
			c.Status(http.StatusOK)
		}
	}
}

func absoluteLink(path string) string {
	firstRune, _ := utf8.DecodeRuneInString(path)
	if string(firstRune) == "/" {
		return config.Config.BaseURL + path
	} else {
		return config.Config.BaseURL + "/" + path
	}
}
