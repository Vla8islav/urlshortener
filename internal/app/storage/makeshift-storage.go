package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
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

func loadDataFromFile(filename string, s Storage) error {
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
		s.AddURLPairInMemory(configuration.ReadFlags().ShortenerBaseURL+"/"+data.ShortURL, data.OriginalURL, data.UUID)

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

func NewMakeshiftStorage() (Storage, error) {
	instance := new(MakeshiftStorage)
	instance.urlToShort = make(map[string]string)
	instance.shortToURL = make(map[string]string)
	instance.uuidList = make(map[string]struct{})
	instance.filePath = configuration.ReadFlags().FileStoragePath

	err := loadDataFromFile(instance.filePath, instance)
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

func (s MakeshiftStorage) AddURLPair(shortenedURL string, fullURL string, uuidStr string) {
	if _, found := s.uuidList[uuidStr]; found {
		return
	}
	s.AddURLPairInMemory(shortenedURL, fullURL, uuidStr)
	writeIntoFile(s.filePath, dataStorageRecord{UUID: uuidStr,
		ShortURL: strings.TrimPrefix(
			strings.TrimPrefix(shortenedURL, configuration.ReadFlags().ShortenerBaseURL), "/"), OriginalURL: fullURL})
}

func (s MakeshiftStorage) AddURLPairInMemory(shortenedURL string, fullURL string, uuidStr string) {
	mu.Lock()
	defer mu.Unlock()
	s.urlToShort[fullURL] = shortenedURL
	s.shortToURL[shortenedURL] = fullURL
	s.uuidList[uuidStr] = struct{}{}
}

func (s MakeshiftStorage) GetFullURL(shortenedURL string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()
	value, exists := s.shortToURL[shortenedURL]
	return value, exists
}

func (s MakeshiftStorage) GetShortenedURL(fullURL string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()
	value, exists := s.urlToShort[fullURL]
	return value, exists
}

func (s MakeshiftStorage) Ping() error {
	mu.Lock()
	defer mu.Unlock()
	f, err := openFileForReading(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
