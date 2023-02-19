package storage

import (
	"hash/fnv"
	"strconv"

	"github.com/webmstk/shorter/internal/config"
)

type Storage interface {
	SaveLongURL(longURL string) (shortURL string, err error)
	GetLongURL(shortURL string) (longURL string, ok bool)
}

func NewStorage() Storage {
	if config.Config.FileStoragePath == "" {
		return &StorageMap{data: make(map[string]string)}
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
