package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"os"
	"strings"
	"sync"
)

var mu sync.Mutex

type dataStorageRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func openFileForReading(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
}

func loadDataFromFile(ctx context.Context, filename string, s Storage) error {
	f, err := openFileForReading(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fileScanner := bufio.NewScanner(f)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		t := fileScanner.Text()
		fmt.Println(t)
		var data dataStorageRecord
		err = json.Unmarshal([]byte(t), &data)
		if err != nil {
			return err
		}
		fmt.Println(data)
		s.AddURLPairInMemory(ctx, configuration.ReadFlags().ShortenerBaseURL+"/"+data.ShortURL, data.OriginalURL, data.UUID, auth.DefaultUserID)

	}
	return nil
}

func writeIntoFile(filename string, data dataStorageRecord) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	dataString, err := json.Marshal(data)
	if err != nil {
		return err
	}
	dataString = append(dataString, '\n')
	_, err = f.Write(dataString)
	return err

}

func NewMakeshiftStorage(ctx context.Context) (Storage, error) {
	instance := new(MakeshiftStorage)
	instance.urlToShort = make(map[string]string)
	instance.shortToURL = make(map[string]string)
	instance.uuidList = make(map[string]struct{})
	instance.filePath = configuration.ReadFlags().FileStoragePath

	err := loadDataFromFile(ctx, instance.filePath, instance)
	if err != nil {
		return instance, err
	}
	return instance, nil
}

type MakeshiftStorage struct {
	urlToShort map[string]string
	shortToURL map[string]string
	uuidList   map[string]struct{}
	filePath   string
}

func (s MakeshiftStorage) Close() {}

func (s MakeshiftStorage) AddURLPair(ctx context.Context, shortenedURL string, fullURL string, uuidStr string, userID int) {
	if _, found := s.uuidList[uuidStr]; found {
		return
	}
	s.AddURLPairInMemory(ctx, shortenedURL, fullURL, uuidStr, userID)
	writeIntoFile(s.filePath, dataStorageRecord{UUID: uuidStr,
		ShortURL: strings.TrimPrefix(
			strings.TrimPrefix(shortenedURL, configuration.ReadFlags().ShortenerBaseURL), "/"), OriginalURL: fullURL})
}

func (s MakeshiftStorage) AddURLPairInMemory(ctx context.Context, shortenedURL string, fullURL string, uuidStr string, userID int) {
	mu.Lock()
	defer mu.Unlock()
	s.urlToShort[fullURL] = shortenedURL
	s.shortToURL[shortenedURL] = fullURL
	s.uuidList[uuidStr] = struct{}{}
}

func (s MakeshiftStorage) GetFullURL(ctx context.Context, shortenedURL string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()
	value, exists := s.shortToURL[shortenedURL]
	return value, exists
}

func (s MakeshiftStorage) GetShortenedURL(ctx context.Context, fullURL string) (string, int, bool) {
	mu.Lock()
	defer mu.Unlock()
	value, exists := s.urlToShort[fullURL]
	return value, auth.DefaultUserID, exists
}

func (s MakeshiftStorage) Ping(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()
	f, err := openFileForReading(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func (s MakeshiftStorage) GetAllURLRecordsByUser(ctx context.Context, userID int) ([]URLPair, error) {
	return []URLPair{}, nil // TODO: actually implement
}

func (s MakeshiftStorage) GetNewUserID(ctx context.Context) (int, error) {
	return auth.DefaultUserID, nil
}

func (s MakeshiftStorage) DeleteURL(ctx context.Context, shortenedURL string) error {
	return nil // TODO: actually implement
}
