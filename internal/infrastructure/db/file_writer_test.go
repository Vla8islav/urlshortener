package db

import (
	"github.com/stretchr/testify/assert"
	"math/rand/v2"
	"path/filepath"
	"strconv"
	"testing"
)

func tempFilePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, strconv.Itoa(rand.Int())+"sample_db.json")
}

func randomRecord() Record {
	return Record{
		ID:           rand.Int64(),
		FullURL:      strconv.Itoa(rand.Int()),
		ShortenedURL: strconv.Itoa(rand.Int()),
	}
}

func TestDataWriterReader_WriteTwiceReadOnce(t *testing.T) {
	t.Parallel()

	path := tempFilePath(t)

	w, err := NewDataWriter(path)
	assert.NoError(t, err, "NewDataWriter error")

	recs1 := []Record{randomRecord(), randomRecord()}
	err = w.WriteRecords(&recs1)
	assert.NoError(t, err, "WriteRecords error")

	recs2 := []Record{randomRecord(), randomRecord()}
	err = w.WriteRecords(&recs2)
	assert.NoError(t, err, "WriteRecords second time error")

	err = w.Close()
	assert.NoError(t, err, "writer Close error: %v", err)

	r, err := NewDataReader(path)
	assert.NoError(t, err, "NewDataReader error")
	defer r.Close()

	got, err := r.ReadRecords()
	assert.NoError(t, err, "ReadRecords error")

	assert.Len(t, *got, 2, "ReadRecords length")
	assert.Equal(t, recs2, *got, "ReadRecords content")
}

func TestNewDataReader_FileNotFound(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), strconv.Itoa(rand.Int()))
	_, err := NewDataReader(path)
	assert.Error(t, err, "NewDataReader should return error for missing file")
}
