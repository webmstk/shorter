package storage

import (
	"hash/fnv"
	"strconv"
)

type Storage interface {
	SaveLongURL(longURL string) (shortURL string, err error)
	GetLongURL(shortURL string) (longURL string, ok bool)
}

type MapStorage map[string]string

func NewStorage() MapStorage {
	return make(MapStorage)
}

func (urls MapStorage) SaveLongURL(longURL string) (shortURL string, err error) {
	shortURL, err = GenerateShortLink(longURL)
	if err != nil {
		return "", err
	}
	// не буду искать, была ли сохранена ссылка ранее, перезапишу ключ, благо хеш
	// сгенерится такой же
	urls[shortURL] = longURL
	return shortURL, nil
}

func (urls MapStorage) GetLongURL(shortURL string) (longURL string, ok bool) {
	longURL, ok = urls[shortURL]
	return
}

func GenerateShortLink(s string) (string, error) {
	h := fnv.New32()
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(h.Sum32())), nil
}
