package storage

import (
	"hash/fnv"
	"strconv"
	"sync"
)

type Storage interface {
	SaveLongURL(longURL string) (shortURL string, err error)
	GetLongURL(shortURL string) (longURL string, ok bool)
}

type MapStorage struct {
	mu   sync.Mutex
	data map[string]string
}

func NewStorage() *MapStorage {
	return &MapStorage{data: make(map[string]string)}
}

func (urls *MapStorage) SaveLongURL(longURL string) (shortURL string, err error) {
	shortURL, err = GenerateShortLink(longURL)
	if err != nil {
		return "", err
	}
	// не буду искать, была ли сохранена ссылка ранее, перезапишу ключ, благо хеш
	// сгенерится такой же
	urls.mu.Lock()
	defer urls.mu.Unlock()
	urls.data[shortURL] = longURL
	return shortURL, nil
}

func (urls *MapStorage) GetLongURL(shortURL string) (longURL string, ok bool) {
	urls.mu.Lock()
	defer urls.mu.Unlock()
	longURL, ok = urls.data[shortURL]
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
