package db

import (
	"github.com/Vla8islav/urlshortener/internal/infrastructure/errlist"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func TestGetInstance_Singleton(t *testing.T) {
	a := GetInstance()
	b := GetInstance()

	if a != b {
		t.Fatalf("expected GetInstance to return the same pointer, got %p and %p", a, b)
	}
}

func randomURL() string {
	return "https://example.com/" + strconv.Itoa(rand.Int())
}

func TestGetByFull_CreatesAndReturnsSameShortURL(t *testing.T) {
	db := GetInstance()

	full := randomURL()

	first, err := db.GetByFull(full)
	assert.NoError(t, err, "unexpected error on first GetByFull")

	assert.Equal(t, full, first.FullURL, "FullURL mismatch on first GetByFull")

	assert.NotEqual(t, first.ShortenedURL, "", "ShortenedURL shouldn't be empty")

	second, err := db.GetByFull(full)
	assert.NoError(t, err, "unexpected error on second GetByFull")

	assert.Equal(t, first.ShortenedURL, second.ShortenedURL,
		"received different ShortenedURL on two requests")

}

func TestGetByShortened_ReturnsFullURL(t *testing.T) {
	db := GetInstance()

	full := randomURL()
	byFull, err := db.GetByFull(full)
	assert.NoError(t, err, "unexpected error on GetByFull")

	byShortened, err := db.GetByShortened(byFull.ShortenedURL)
	assert.NoError(t, err, "unexpected error on GetByShortened")
	assert.Equal(t, full, byShortened.FullURL, "FullURL mismatch on GetByShortened")
	assert.Equal(t, byFull.ShortenedURL, byShortened.ShortenedURL, "ShortenedURL mismatch on GetByShortened")
}

func TestGetByShortened_UnknownURL(t *testing.T) {
	db := GetInstance()

	_, err := db.GetByShortened("nonexistent123")
	assert.Error(t, err, "expected error for nonexistent short URL")
	assert.ErrorAs(t, err, &errlist.ErrURLNotFound, "expected ErrURLNotFound type")
}
