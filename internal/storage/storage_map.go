package storage

import (
	"sync"
)

type StorageMap struct {
	mu   sync.Mutex
	data map[string]string
}

func (urls *StorageMap) SaveLongURL(longURL string) (shortURL string, err error) {
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

func (urls *StorageMap) GetLongURL(shortURL string) (longURL string, ok bool) {
	urls.mu.Lock()
	defer urls.mu.Unlock()
	longURL, ok = urls.data[shortURL]
	return
}
