package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/storage"
)

func UserCookie(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Cookie("user_id")
		userToken, _ := c.Cookie("user_token")
		if userID == "" {
			userID = store.CreateUser(context.Background())
			userToken = SignCookie(userID)
		} else if !IsTokenValid(userID, userToken) {
			userID = store.CreateUser(context.Background())
			userToken = SignCookie(userID)
		}

		c.Set("user_id", userID)

		host := strings.Split(config.Config.ServerAddress, ":")[0]
		c.SetCookie("user_id", userID, config.Config.CookieTTLSeconds, "/", host, false, true)
		c.SetCookie("user_token", userToken, config.Config.CookieTTLSeconds, "/", host, false, true)
	}
}

func SignCookie(content string) string {
	key := []uint8(config.Config.CookieSalt)

	h := hmac.New(sha256.New, key)
	h.Write([]byte(content))
	dst := h.Sum(nil)

	return fmt.Sprintf("%x", dst)
}

func IsTokenValid(userID string, userToken string) bool {
	return SignCookie(userID) == userToken
}
