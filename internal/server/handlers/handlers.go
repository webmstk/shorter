package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"

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

		userID, _ := c.Cookie("user_id")
		userToken, _ := c.Cookie("user_token")
		if userID == "" {
			userID = storage.CreateUser()
			userToken = signCookie(userID)
		} else if !isTokenValid(userID, userToken) {
			userID = storage.CreateUser()
			userToken = signCookie(userID)
		}

		shortURL, err := storage.SaveLongURL(longURL, userID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to generate short link")
			return
		}

		c.SetCookie("user_id", userID, config.Config.CookieTTLSeconds, "/", config.Config.ServerAddress, false, true)
		c.SetCookie("user_token", userToken, config.Config.CookieTTLSeconds, "/", config.Config.ServerAddress, false, true)
		c.String(http.StatusCreated, absoluteLink(shortURL))
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

func HandlerPing(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := pgx.Connect(context.Background(), config.Config.DatabaseDSN)
		if err != nil {
			log.Print("DB failure: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		defer conn.Close(context.Background())

		_, err = conn.Query(context.Background(), "SELECT 1")
		if err != nil {
			log.Print("DB failure: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
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
