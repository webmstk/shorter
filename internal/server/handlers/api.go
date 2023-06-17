package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/server/util"
	"github.com/webmstk/shorter/internal/storage"
)

func HandlerAPIShorten(storage storage.Storage) gin.HandlerFunc {
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

		shortURL, err := storage.SaveLongURL(reqBody.URL, "")
		if err != nil {
			response := util.WriteJSONError(c, "internal error")
			c.String(http.StatusInternalServerError, response)
			return
		}
		respBody := struct {
			Result string `json:"result"`
		}{
			Result: absoluteLink(shortURL),
		}

		c.JSON(http.StatusCreated, respBody)
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

		output, err := store.SaveBatch(reqBody)
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
		if !isTokenValid(userID, userToken) {
			c.Status(http.StatusNoContent)
			return
		}

		links, ok := storage.GetUserLinks(userID)
		if !ok {
			c.Status(http.StatusNoContent)
			return
		}

		var response []map[string]string
		for _, link := range links {
			longURL, ok := storage.GetLongURL(link)
			if ok {
				fields := map[string]string{"short_url": absoluteLink(link), "original_url": longURL}
				response = append(response, fields)
			}
		}
		c.JSON(http.StatusOK, response)
	}
}
