package storage

import (
	"encoding/json"
	"os"
	"sync"
)

type StorageFile struct {
	mu       sync.Mutex
	filePath string
}

func (fs *StorageFile) SaveLongURL(longURL string) (shortURL string, err error) {
	shortURL, err = GenerateShortLink(longURL)
	if err != nil {
		return "", err
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := os.OpenFile(fs.filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	storage, err := parseFile(file)
	if err != nil {
		return "", err
	}

	storage[shortURL] = longURL
	data, err := json.MarshalIndent(&storage, "", "  ")
	if err != nil {
		return "", err
	}

	file.WriteAt(data, 0)

	return shortURL, nil
}

func (fs *StorageFile) GetLongURL(shortURL string) (longURL string, ok bool) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := os.OpenFile(fs.filePath, os.O_RDONLY, 0644)
	if err != nil {
		return "", false
	}
	defer file.Close()

	storage, err := parseFile(file)
	if err != nil {
		return "", false
	}

	longURL, ok = storage[shortURL]
	return
}

func parseFile(file *os.File) (content map[string]string, err error) {
	stat, err := file.Stat()
	if err != nil {
		return content, err
	}

	buf := make([]byte, stat.Size())
	file.Read(buf)

	if len(buf) == 0 {
		return map[string]string{}, nil
	}

	err = json.Unmarshal(buf, &content)
	if err != nil {
		return content, err
	}

	return content, nil
}
