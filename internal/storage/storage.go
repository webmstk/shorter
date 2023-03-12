package storage

import (
	"hash/fnv"
	"strconv"

	"github.com/webmstk/shorter/internal/config"
)

type Storage interface {
	SaveLongURL(longURL, userID string) (shortURL string, err error)
	GetLongURL(shortURL string) (longURL string, ok bool)
	CreateUser() string
	GetUserLinks(UserID string) (links []string, ok bool)
}

type table map[string]any

func NewStorage() Storage {
	if config.Config.FileStoragePath == "" {
		return &StorageMap{data: make(map[string]table)}
	} else {
		return &StorageFile{filePath: config.Config.FileStoragePath}
	}
}

func GenerateShortLink(s string) (string, error) {
	h := fnv.New32()
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(h.Sum32())), nil
}

func getTable(data map[string]table, tableName string) table {
	if data[tableName] == nil {
		data[tableName] = table{}
	}

	return data[tableName]
}
