package db

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Vla8islav/urlshortener/internal/infrastructure/errlist"
)

func resetFileDB() {
	fileDBInstance = nil
	once = sync.Once{}
}

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), strconv.Itoa(rand.Int()))
}

func Test_FileDB_Instance_InitAndSmoke(t *testing.T) {
	defer resetFileDB()

	filename := tempPath(t)

	db := GetFileDBInstance(filename)
	require.NotNil(t, db, "GetFileDBInstance should not return nil")
	require.NotNil(t, db.records, "records must be initialized")
	assert.Equal(t, 0, len(*db.records), "new DB should start empty")
}

func Test_GetByFull_CreatesAndIsIdempotent(t *testing.T) {
	defer resetFileDB()

	filename := tempPath(t)
	db := GetFileDBInstance(filename)
	require.NotNil(t, db)

	full := randomURL()

	u1, err := db.GetByFull(full)
	require.NoError(t, err)
	assert.Equal(t, full, u1.FullURL)
	assert.NotEmpty(t, u1.ShortenedURL)

	u2, err := db.GetByFull(full)
	require.NoError(t, err)
	assert.Equal(t, u1.ShortenedURL, u2.ShortenedURL, "same full URL must map to the same short URL")

	require.NotNil(t, db.records)
	assert.Equal(t, 1, len(*db.records))
	assert.Equal(t, int64(1), db.getHighestID(), "first record should have ID 1")
}

func Test_GetByShortened_GetByShortened(t *testing.T) {
	defer resetFileDB()

	filename := tempPath(t)
	db := GetFileDBInstance(filename)
	require.NotNil(t, db)

	full := randomURL()
	u, err := db.GetByFull(full)
	require.NoError(t, err)
	require.NotEmpty(t, u.ShortenedURL)

	back, err := db.GetByShortened(u.ShortenedURL)
	require.NoError(t, err)
	assert.Equal(t, full, back.FullURL)
	assert.Equal(t, u.ShortenedURL, back.ShortenedURL)
}

func Test_GetByShortened_UnknownReturnsErrURLNotFound(t *testing.T) {
	defer resetFileDB()

	filename := tempPath(t)
	db := GetFileDBInstance(filename)
	require.NotNil(t, db)

	_, err := db.GetByShortened("does-not-exist")
	require.Error(t, err)
	assert.Equal(t, errlist.ErrURLNotFound, err)
}

func Test_PersistenceAcrossInstances_SameFile(t *testing.T) {
	defer resetFileDB()

	filename := tempPath(t)

	db1 := GetFileDBInstance(filename)
	require.NotNil(t, db1)

	url1 := randomURL()
	url2 := randomURL()

	u1, err := db1.GetByFull(url1)
	require.NoError(t, err)
	u2, err := db1.GetByFull(url2)
	require.NoError(t, err)
	require.NotEqual(t, u1.ShortenedURL, u2.ShortenedURL)
	assert.Equal(t, int64(2), db1.getHighestID(), "highest ID should be 2 after two inserts")

	resetFileDB()
	db2 := GetFileDBInstance(filename)
	require.NotNil(t, db2)
	require.NotNil(t, db2.records)

	// Records are reloaded
	assert.Equal(t, 2, len(*db2.records), "should reload two records from file")

	// We can get data back
	back1, err := db2.GetByShortened(u1.ShortenedURL)
	require.NoError(t, err)
	assert.Equal(t, url1, back1.FullURL)

	back2, err := db2.GetByShortened(u2.ShortenedURL)
	require.NoError(t, err)
	assert.Equal(t, url2, back2.FullURL)

	// highest ID stays correct
	assert.Equal(t, int64(2), db2.getHighestID())
}
