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
