package storage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/webmstk/shorter/internal/config"
)

func TestGenerateShortLink(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "value", value: "aaa", want: "1184641920"},
		{name: "same value", value: "aaa", want: "1184641920"},
		{name: "another value", value: "bbb", want: "496612573"},
	}

	for _, tt := range tests {
		actual, _ := GenerateShortLink(tt.value)
		assert.Equal(t, tt.want, actual)
	}
}

func TestSaveLongURL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		user  bool
		want  string
	}{
		{name: "value", value: "https://ya1.ru", want: "3144338548"},
		{name: "same value", value: "https://ya1.ru", want: "3144338548"},
		{name: "another value", value: "https://yandex2.ru", want: "119677240"},
		{name: "value", value: "https://ya5.ru", user: true, want: "269264536"},
	}

	storageMap := &StorageMap{data: make(map[string]table)}
	userID := storageMap.CreateUser(context.Background())

	for _, tt := range tests {
		user := ""
		if tt.user {
			user = userID
		}
		shortURL, _ := storageMap.SaveLongURL(context.Background(), tt.value, user)
		longURL, err := storageMap.GetLongURL(context.Background(), shortURL)

		require.NoError(t, err)
		assert.Equal(t, tt.want, shortURL)
		assert.Equal(t, tt.value, longURL)

		userLinks, err := storageMap.GetUserLinks(context.Background(), user)
		if user != "" {
			require.NoError(t, err)
			assert.Equal(t, tt.want, userLinks[0])
		}
	}

	path := "../../storage/storage_test.json"
	storageFile := &StorageFile{filePath: path}
	userID = storageFile.CreateUser(context.Background())
	defer os.Remove(path)

	for _, tt := range tests {
		user := ""
		if tt.user {
			user = userID
		}
		shortURL, _ := storageFile.SaveLongURL(context.Background(), tt.value, user)
		longURL, err := storageFile.GetLongURL(context.Background(), shortURL)

		require.NoError(t, err)
		assert.Equal(t, tt.want, shortURL)
		assert.Equal(t, tt.value, longURL)

		if user != "" {
			userLinks, err := storageFile.GetUserLinks(context.Background(), user)
			require.NoError(t, err)
			assert.Equal(t, tt.want, userLinks[0])
		}
	}

	if config.Config.DatabaseDSN != "" {
		storageDB, _ := NewStorageDB()
		userID = storageDB.CreateUser(context.Background())
		for _, tt := range tests {
			user := ""
			if tt.user {
				user = userID
			}
			shortURL, _ := storageDB.SaveLongURL(context.Background(), tt.value, user)
			longURL, err := storageDB.GetLongURL(context.Background(), shortURL)

			require.NoError(t, err)
			assert.Equal(t, tt.want, shortURL)
			assert.Equal(t, tt.value, longURL)

			userLinks, err := storageDB.GetUserLinks(context.Background(), user)
			if user != "" {
				require.NoError(t, err)
				assert.Equal(t, tt.want, userLinks[0])
			}
		}
	}
}

func TestGetLongURL(t *testing.T) {
	storageMap := &StorageMap{data: make(map[string]table)}
	link := "https://ya.ru"
	shortURL, _ := storageMap.SaveLongURL(context.Background(), link, "")

	type want struct {
		link string
		ok   bool
	}

	tests := []struct {
		name  string
		value string
		want  want
	}{
		{
			name:  "link exists",
			value: shortURL,
			want: want{
				link: link,
				ok:   true,
			},
		},
		{
			name:  "link does not exist",
			value: "404",
			want: want{
				link: "",
				ok:   false,
			},
		},
	}

	for _, tt := range tests {
		longURL, err := storageMap.GetLongURL(context.Background(), tt.value)
		if tt.want.ok {
			require.NoError(t, err)
		}
		assert.Equal(t, tt.want.link, longURL)
	}

	path := "../../storage/storage_test.json"
	storageFile := &StorageFile{filePath: path}
	defer os.Remove(path)

	_, err := storageFile.SaveLongURL(context.Background(), link, "")
	if err != nil {
		panic(err)
	}

	for _, tt := range tests {
		longURL, err := storageFile.GetLongURL(context.Background(), tt.value)
		if tt.want.ok {
			require.NoError(t, err)
		}
		assert.Equal(t, tt.want.link, longURL)
	}

	if config.Config.DatabaseDSN != "" {
		storageDB, _ := NewStorageDB()

		_, err = storageDB.SaveLongURL(context.Background(), link, "")
		if err != nil {
			panic(err)
		}

		for _, tt := range tests {
			longURL, err := storageDB.GetLongURL(context.Background(), tt.value)
			if tt.want.ok {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want.link, longURL)
		}
	}
}

func TestSaveBatch(t *testing.T) {
	storageMap := &StorageMap{data: make(map[string]table)}

	tests := []struct {
		name  string
		input []BatchInput
		want  []BatchOutput
	}{
		{
			name: "test_1",
			input: []BatchInput{
				{CorrelationID: "111", OriginalURL: "http://lelik.ru"},
				{CorrelationID: "222", OriginalURL: "http://bolik.ru"},
			},
			want: []BatchOutput{
				{CorrelationID: "111", ShortURL: "3799407019"},
				{CorrelationID: "222", ShortURL: "2114288767"},
			},
		},
	}

	for _, tt := range tests {
		output, err := storageMap.SaveBatch(context.Background(), tt.input)
		require.NoError(t, err)
		assert.Equal(t, tt.want, output)
	}

	path := "../../storage/storage_test.json"
	storageFile := &StorageFile{filePath: path}
	defer os.Remove(path)

	for _, tt := range tests {
		output, err := storageFile.SaveBatch(context.Background(), tt.input)
		require.NoError(t, err)
		assert.Equal(t, tt.want, output)
	}

	if config.Config.DatabaseDSN != "" {
		storageDB, _ := NewStorageDB()

		for _, tt := range tests {
			output, err := storageDB.SaveBatch(context.Background(), tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, output)
		}
	}
}
