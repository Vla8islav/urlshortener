package db

import (
	"errors"
	"github.com/Vla8islav/urlshortener/internal/domain"
	"github.com/Vla8islav/urlshortener/internal/helpers"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/errlist"
	"sync"
)

func convertToURL(r Record) domain.URL {
	return domain.URL{
		FullURL:      r.FullURL,
		ShortenedURL: r.ShortenedURL,
	}
}

type Record struct {
	ID           int64
	FullURL      string
	ShortenedURL string
}

type MakeshiftFileDB struct {
	records []Record
	writer  *DataWriter
	reader  *DataReader
}

var fileDBInstance *MakeshiftFileDB = nil

func GetFileDBInstance(filename string) *MakeshiftFileDB {
	sync.OnceFunc(func() {
		fileDBInstance = new(MakeshiftFileDB)
		writer, err := NewDataWriter(filename)
		if err != nil {
			return
		}
		reader, err := NewDataReader(filename)
		if err != nil {
			return
		}
		fileDBInstance.writer = writer
		fileDBInstance.reader = reader
	})()

	return fileDBInstance
}

func (s MakeshiftFileDB) GetByFull(fullURL string) (domain.URL, error) {

	searchResult, err := s.search(fullURL, "")
	if err == nil {
		return convertToURL(searchResult), nil
	}

	if errors.Is(err, errlist.ErrURLNotFound) {
		newRecord := Record{
			ID:           s.getHighestID() + 1,
			FullURL:      fullURL,
			ShortenedURL: helpers.GenerateShortenedURLUID(),
		}
		s.records = append(s.records, newRecord)
		err = s.writer.WriteRecords(s.records)
		if err != nil {
			return domain.URL{}, err
		}

	}

	return domain.URL{}, err
}

func (s MakeshiftFileDB) GetByShortened(shortURL string) (domain.URL, error) {
	searchResult, err := s.search("", shortURL)
	if err == nil {
		return convertToURL(searchResult), nil
	}

	return domain.URL{}, err
}

func (s MakeshiftFileDB) search(fullURL, shortenedURL string) (Record, error) {
	for _, record := range s.records {
		if record.FullURL == fullURL && len(fullURL) > 0 {
			return record, nil
		}
		if record.ShortenedURL == shortenedURL && len(shortenedURL) > 0 {
			return record, nil
		}
	}
	return Record{}, errlist.ErrURLNotFound
}

func (s MakeshiftFileDB) getHighestID() int64 {
	id := int64(1)
	for _, record := range s.records {
		if record.ID > id {
			id = record.ID + 1
		}
	}
	return id
}
