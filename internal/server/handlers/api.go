package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/server/middlewares"
	"github.com/webmstk/shorter/internal/server/util"
	"github.com/webmstk/shorter/internal/storage"
)

func HandlerAPIShorten(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		if !util.ValidateContentTypeJSON(c, "application/json") {
			return
		}
		body, err := util.ReadBody(c)
		reqBody := struct {
			URL string `json:"url"`
		}{}
		if err != nil {
			response := util.WriteJSONError(c, "url is not valid")
			c.String(http.StatusBadRequest, response)
			return
		}

		err = json.Unmarshal([]byte(body), &reqBody)
		if err != nil || reqBody.URL == "" {
			response := util.WriteJSONError(c, "url is not valid")
			c.String(http.StatusBadRequest, response)
			return
		}

		status := http.StatusCreated
		shortURL, err := store.SaveLongURL(c, reqBody.URL, "")
		if err != nil {
			var linkExistError *storage.LinkExistError
			if errors.As(err, &linkExistError) {
				status = http.StatusConflict
			} else {
				response := util.WriteJSONError(c, "internal error")
				c.String(http.StatusInternalServerError, response)
				return
			}
		}
		respBody := struct {
			Result string `json:"result"`
		}{
			Result: absoluteLink(shortURL),
		}

		c.JSON(status, respBody)
	}
}

func HandlerAPIShortenBatch(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		if !util.ValidateContentTypeJSON(c, "application/json") {
			return
		}

		body, _ := util.ReadBody(c)
		var reqBody []storage.BatchInput

		err := json.Unmarshal([]byte(body), &reqBody)
		if err != nil {
			response := util.WriteJSONError(c, "url is not valid")
			c.String(http.StatusBadRequest, response)
			return
		}

		output, err := store.SaveBatch(c, reqBody)
		if err != nil {
			response := util.WriteJSONError(c, "failed to save some links")
			c.String(http.StatusBadRequest, response)
			return
		}

		var result []storage.BatchOutput
		for _, elem := range output {
			elem.ShortURL = absoluteLink(elem.ShortURL)
			result = append(result, elem)
		}
		c.JSON(http.StatusCreated, result)
	}
}

func HandlerAPIUserUrls(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Cookie("user_id")
		userToken, _ := c.Cookie("user_token")

		var links []string
		if !middlewares.IsTokenValid(userID, userToken) {
			c.Status(http.StatusNoContent)
			return
		}

		links, ok := storage.GetUserLinks(c, userID)
		if !ok {
			c.Status(http.StatusNoContent)
			return
		}

		var response []map[string]string
		for _, link := range links {
			longURL, ok := storage.GetLongURL(c, link)
			if ok {
				fields := map[string]string{"short_url": absoluteLink(link), "original_url": longURL}
				response = append(response, fields)
			}
		}
		c.JSON(http.StatusOK, response)
	}
}
