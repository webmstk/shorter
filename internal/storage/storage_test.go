package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
	urls := NewStorage()

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "value", value: "https://ya.ru", want: "2138586483"},
		{name: "same value", value: "https://ya.ru", want: "2138586483"},
		{name: "another value", value: "https://yandex.ru", want: "1250700976"},
	}

	for _, tt := range tests {
		shortURL, err := urls.SaveLongURL(tt.value)
		longURL, ok := urls.GetLongURL(shortURL)

		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, tt.want, shortURL)
		assert.Equal(t, tt.value, longURL)
	}
}

func TestGetLongURL(t *testing.T) {
	urls := NewStorage()
	link := "https://ya.ru"
	shortURL, _ := urls.SaveLongURL(link)

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
		longURL, ok := urls.GetLongURL(tt.value)
		assert.Equal(t, tt.want.ok, ok)
		assert.Equal(t, tt.want.link, longURL)
	}
}
