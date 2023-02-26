package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/config"
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

		shortURL, err := storage.SaveLongURL(reqBody.URL)
		if err != nil {
			response := util.WriteJSONError(c, "internal error")
			c.String(http.StatusInternalServerError, response)
			return
		}
		respBody := struct {
			Result string `json:"result"`
		}{
			Result: config.Config.BaseURL + "/" + shortURL,
		}

		result, _ := json.Marshal(respBody)
		c.String(http.StatusCreated, string(result))
	}
}
