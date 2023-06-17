package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
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

		userID, _ := c.Cookie("user_id")
		userToken, _ := c.Cookie("user_token")
		if userID == "" {
			userID = store.CreateUser()
			userToken = signCookie(userID)
		} else if !isTokenValid(userID, userToken) {
			userID = store.CreateUser()
			userToken = signCookie(userID)
		}

		status := http.StatusCreated
		shortURL, err := store.SaveLongURL(longURL, userID)
		fmt.Println(err)
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

func isTokenValid(userID string, userToken string) bool {
	return signCookie(userID) == userToken
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

func HandlerPing(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch s := store.(type) {
		case *storage.StorageDB:
			err := s.Ping()
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

func signCookie(content string) string {
	key := []uint8(config.Config.CookieSalt)

	h := hmac.New(sha256.New, key)
	h.Write([]byte(content))
	dst := h.Sum(nil)

	return fmt.Sprintf("%x", dst)
}

func absoluteLink(path string) string {
	firstRune, _ := utf8.DecodeRuneInString(path)
	if string(firstRune) == "/" {
		return config.Config.BaseURL + path
	} else {
		return config.Config.BaseURL + "/" + path
	}
}
