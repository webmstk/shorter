package storage

import (
	"sync"

	"github.com/google/uuid"
)

type StorageMap struct {
	mu   sync.Mutex
	data map[string]table
}

func (storage *StorageMap) SaveLongURL(longURL, userID string) (shortURL string, err error) {
	shortURL, err = GenerateShortLink(longURL)
	if err != nil {
		return "", err
	}
	// не буду искать, была ли сохранена ссылка ранее, перезапишу ключ, благо хеш
	// сгенерится такой же
	storage.mu.Lock()
	defer storage.mu.Unlock()

	getTable(storage.data, "links")[shortURL] = longURL

	if userID != "" {
		userLinks := getTable(storage.data, "user_links")

		if userLinks[userID] == nil {
			userLinks[userID] = []string{shortURL}
		} else {
			userLinks[userID] = append(userLinks[userID].([]string), shortURL)
		}
	}
	return shortURL, nil
}

func (storage *StorageMap) GetLongURL(shortURL string) (longURL string, ok bool) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	value, ok := getTable(storage.data, "links")[shortURL]
	if ok {
		longURL = value.(string)
	}
	return
}

func (storage *StorageMap) GetUserLinks(userID string) (links []string, ok bool) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	value, ok := getTable(storage.data, "user_links")[userID]
	if ok {
		links = value.([]string)
	}
	return
}

func (storage *StorageMap) CreateUser() string {
	return uuid.New().String()
}
