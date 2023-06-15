package storage

import (
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
		{name: "value", value: "https://ya1.ru", user: true, want: "3144338548"},
	}

	storageMap := &StorageMap{data: make(map[string]table)}
	userID := storageMap.CreateUser()

	for _, tt := range tests {
		user := ""
		if tt.user {
			user = userID
		}
		shortURL, err := storageMap.SaveLongURL(tt.value, user)
		longURL, ok := storageMap.GetLongURL(shortURL)

		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, tt.want, shortURL)
		assert.Equal(t, tt.value, longURL)

		userLinks, ok := storageMap.GetUserLinks(user)
		if user != "" {
			require.True(t, ok)
			assert.Equal(t, tt.want, userLinks[0])
		}
	}

	path := "../../storage/storage_test.json"
	storageFile := &StorageFile{filePath: path}
	userID = storageFile.CreateUser()
	defer os.Remove(path)

	for _, tt := range tests {
		user := ""
		if tt.user {
			user = userID
		}
		shortURL, err := storageFile.SaveLongURL(tt.value, user)
		longURL, ok := storageFile.GetLongURL(shortURL)

		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, tt.want, shortURL)
		assert.Equal(t, tt.value, longURL)

		if user != "" {
			userLinks, ok := storageFile.GetUserLinks(user)
			require.True(t, ok)
			assert.Equal(t, tt.want, userLinks[0])
		}
	}

	if config.Config.DatabaseDSN != "" {
		storageDB := NewStorageDB()
		userID = storageDB.CreateUser()
		for _, tt := range tests {
			user := ""
			if tt.user {
				user = userID
			}
			shortURL, err := storageDB.SaveLongURL(tt.value, user)
			longURL, ok := storageDB.GetLongURL(shortURL)

			require.NoError(t, err)
			require.True(t, ok)
			assert.Equal(t, tt.want, shortURL)
			assert.Equal(t, tt.value, longURL)

			userLinks, ok := storageDB.GetUserLinks(user)
			if user != "" {
				require.True(t, ok)
				assert.Equal(t, tt.want, userLinks[0])
			}
		}
	}
}

func TestGetLongURL(t *testing.T) {
	storageMap := &StorageMap{data: make(map[string]table)}
	link := "https://ya.ru"
	shortURL, _ := storageMap.SaveLongURL(link, "")

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
		longURL, ok := storageMap.GetLongURL(tt.value)
		assert.Equal(t, tt.want.ok, ok)
		assert.Equal(t, tt.want.link, longURL)
	}

	path := "../../storage/storage_test.json"
	storageFile := &StorageFile{filePath: path}
	defer os.Remove(path)

	_, err := storageFile.SaveLongURL(link, "")
	if err != nil {
		panic(err)
	}

	for _, tt := range tests {
		longURL, ok := storageFile.GetLongURL(tt.value)
		assert.Equal(t, tt.want.ok, ok)
		assert.Equal(t, tt.want.link, longURL)
	}

	if config.Config.DatabaseDSN != "" {
		storageDB := NewStorageDB()

		_, err = storageDB.SaveLongURL(link, "")
		if err != nil {
			panic(err)
		}

		for _, tt := range tests {
			longURL, ok := storageDB.GetLongURL(tt.value)
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.link, longURL)
		}
	}
}
