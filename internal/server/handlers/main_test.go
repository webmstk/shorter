package handlers

import (
	"os"
	"testing"

	"github.com/webmstk/shorter/internal/tests"
)

func TestMain(m *testing.M) {
	tests.Setup()
	os.Exit(m.Run())
}
