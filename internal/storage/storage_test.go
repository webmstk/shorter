package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/webmstk/shorter/internal/tests"
)

func TestMain(m *testing.M) {
	tests.Setup()
	os.Exit(m.Run())
}

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
		name   string
		value  string
		userID string
		want   string
	}{
		{name: "value", value: "https://ya.ru", want: "2138586483"},
		{name: "same value", value: "https://ya.ru", want: "2138586483"},
		{name: "another value", value: "https://yandex.ru", want: "1250700976"},
		{name: "value", value: "https://ya.ru", userID: "123", want: "2138586483"},
	}

	storageMap := &StorageMap{data: make(map[string]table)}

	for _, tt := range tests {
		shortURL, err := storageMap.SaveLongURL(tt.value, tt.userID)
		longURL, ok := storageMap.GetLongURL(shortURL)

		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, tt.want, shortURL)
		assert.Equal(t, tt.value, longURL)

		userLinks, ok := storageMap.GetUserLinks(tt.userID)
		if tt.userID != "" {
			require.True(t, ok)
			assert.Equal(t, tt.want, userLinks[0])
		}
	}

	path := "../../storage/storage_test.json"
	storageFile := &StorageFile{filePath: path}
	defer os.Remove(path)

	for _, tt := range tests {
		shortURL, err := storageFile.SaveLongURL(tt.value, tt.userID)
		longURL, ok := storageFile.GetLongURL(shortURL)

		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, tt.want, shortURL)
		assert.Equal(t, tt.value, longURL)

		if tt.userID != "" {
			userLinks, ok := storageFile.GetUserLinks(tt.userID)
			require.True(t, ok)
			assert.Equal(t, tt.want, userLinks[0])
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
}
