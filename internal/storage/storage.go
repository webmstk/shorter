package storage

import (
	"context"
	"hash/fnv"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webmstk/shorter/internal/config"
)

type Storage interface {
	SaveLongURL(longURL, userID string) (shortURL string, err error)
	GetLongURL(shortURL string) (longURL string, ok bool)
	CreateUser() string
	GetUserLinks(UserID string) (links []string, ok bool)
	SaveBatch(records []BatchInput) ([]BatchOutput, error)
}

type BatchInput struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchOutput struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type LinkExistError struct {
	ShortURL string
}

func (lee *LinkExistError) Error() string {
	return "link exists"
}

func NewLinkExistError(shortURL string) error {
	return &LinkExistError{
		ShortURL: shortURL,
	}
}

type table map[string]any

func NewStorage() Storage {
	if config.Config.DatabaseDSN != "" {
		return NewStorageDB()
	} else if config.Config.FileStoragePath == "" {
		return NewStorageMap()
	} else {
		return NewStorageFile()
	}
}

func NewStorageMap() *StorageMap {
	return &StorageMap{data: make(map[string]table)}
}

func NewStorageFile() *StorageFile {
	return &StorageFile{filePath: config.Config.FileStoragePath}
}

func NewStorageDB() *StorageDB {
	pool, err := pgxpool.New(context.Background(), config.Config.DatabaseDSN)
	if err != nil {
		log.Fatal("DB failure: ", err)
	}

	return &StorageDB{pool: pool}
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
