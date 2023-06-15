package storage

import (
	"os"
	"testing"

	"github.com/webmstk/shorter/internal/tests"
)

func TestMain(m *testing.M) {
	tests.Setup()
	code := m.Run()
	os.Exit(code)
}
