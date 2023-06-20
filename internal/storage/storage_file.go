package storage

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/google/uuid"
)

type StorageFile struct {
	mu       sync.Mutex
	filePath string
}

func (fs *StorageFile) SaveLongURL(_ctx context.Context, longURL, userID string) (shortURL string, err error) {
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

	conflict := false
	table := getTable(storage, "links")
	if _, ok := table[shortURL]; ok {
		conflict = true
	}
	getTable(storage, "links")[shortURL] = longURL

	if userID != "" {
		userLinks := getTable(storage, "user_links")

		if userLinks[userID] == nil {
			userLinks[userID] = []string{shortURL}
		} else {
			userLinks[userID] = append(userLinks[userID].([]string), shortURL)
		}
	}

	data, err := json.MarshalIndent(&storage, "", "  ")
	if err != nil {
		return "", err
	}

	file.WriteAt(data, 0)

	if conflict {
		return shortURL, NewLinkExistError(shortURL)
	}

	return shortURL, nil
}

func (fs *StorageFile) GetLongURL(_ctx context.Context, shortURL string) (longURL string, ok bool) {
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

	value, ok := getTable(storage, "links")[shortURL]
	if ok {
		longURL = value.(string)
	}
	return
}

func (fs *StorageFile) GetUserLinks(_ctx context.Context, userID string) (links []string, ok bool) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := os.OpenFile(fs.filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	storage, err := parseFile(file)
	if err != nil {
		return nil, false
	}

	values, ok := getTable(storage, "user_links")[userID]
	if ok {
		for _, value := range values.([]interface{}) {
			links = append(links, value.(string))
		}
	}
	return
}

func (fs *StorageFile) CreateUser(_ctx context.Context) string {
	return uuid.New().String()
}

func (fs *StorageFile) SaveBatch(ctx context.Context, records []BatchInput) ([]BatchOutput, error) {
	var output []BatchOutput
	for _, record := range records {
		shortURL, err := fs.SaveLongURL(ctx, record.OriginalURL, "")
		if err != nil {
			return output, err
		}

		batchOutput := BatchOutput{
			CorrelationID: record.CorrelationID,
			ShortURL:      shortURL,
		}
		output = append(output, batchOutput)
	}

	return output, nil
}

func parseFile(file *os.File) (content map[string]table, err error) {
	stat, err := file.Stat()
	if err != nil {
		return content, err
	}

	buf := make([]byte, stat.Size())
	file.Read(buf)

	if len(buf) == 0 {
		return map[string]table{}, nil
	}

	err = json.Unmarshal(buf, &content)
	if err != nil {
		return content, err
	}

	return content, nil
}
