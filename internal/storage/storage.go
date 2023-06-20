package storage

import (
	"context"
	"errors"
	"hash/fnv"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webmstk/shorter/internal/config"
)

type Storage interface {
	SaveLongURL(ctx context.Context, longURL, userID string) (shortURL string, err error)
	GetLongURL(ctx context.Context, shortURL string) (longURL string, ok bool)
	CreateUser(ctx context.Context) string
	GetUserLinks(ctx context.Context, UserID string) (links []string, ok bool)
	SaveBatch(ctx context.Context, records []BatchInput) ([]BatchOutput, error)
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

func NewStorage() (Storage, error) {
	if config.Config.DatabaseDSN != "" {
		return NewStorageDB()
	} else if config.Config.FileStoragePath == "" {
		return NewStorageMap(), nil
	} else {
		return NewStorageFile()
	}
}

func NewStorageMap() *StorageMap {
	return &StorageMap{data: make(map[string]table)}
}

func NewStorageFile() (*StorageFile, error) {
	_, err := os.Stat(config.Config.FileStoragePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, err := os.Create(config.Config.FileStoragePath)
			if err != nil {
				return nil, err
			}
			defer f.Close()

		} else {
			return nil, err
		}
	}

	return &StorageFile{filePath: config.Config.FileStoragePath}, nil
}

func NewStorageDB() (*StorageDB, error) {
	pool, err := pgxpool.New(context.Background(), config.Config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	return &StorageDB{pool: pool}, nil
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
