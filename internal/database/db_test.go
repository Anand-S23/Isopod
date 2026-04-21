package database

import (
	"context"
	"testing"
)

func TestInitDB_invalidURL(t *testing.T) {
	_, err := InitDB(context.Background(), "invalid://not-postgres")
	if err == nil {
		t.Fatal("expected error for invalid driver / URL")
	}
}
