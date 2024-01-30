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

type DataStorageRecord struct {
	UUID        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

func loadDataFromFile(filename string, s Storage) error {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(f)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		t := fileScanner.Text()
		fmt.Println(t)
		var data DataStorageRecord
		err = json.Unmarshal([]byte(t), &data)
		if err != nil {
			return err
		}
		fmt.Println(data)
		s.AddUrlPairInMemory(configuration.ReadFlags().ShortenerBaseURL+"/"+data.ShortUrl, data.OriginalUrl, data.UUID)

	}
	return nil
}

func NewMakeshiftStorage() (Storage, error) {
	instance := new(MakeshiftStorage)
	instance.urlToShort = make(map[string]string)
	instance.shortToURL = make(map[string]string)
	instance.uuidList = make(map[string]struct{})
	err := loadDataFromFile(configuration.ReadFlags().FileStoragePath, instance)
	if err != nil {
		return instance, err
	}
	return instance, nil
}

type Storage interface {
	AddURLPair(shortenedURL string, fullURL string, uuidStr string)
	AddUrlPairInMemory(shortenedURL string, fullURL string, uuidStr string)
	GetFullURL(shortenedURL string) (string, bool)
	GetShortenedURL(fullURL string) (string, bool)
}

type MakeshiftStorage struct {
	urlToShort map[string]string
	shortToURL map[string]string
	uuidList   map[string]struct{}
}

func writeIntoFile(filename string, data DataStorageRecord) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	dataString, err := json.Marshal(data)
	if err != nil {
		return err
	}
	dataString = append(dataString, '\n')
	_, err = f.Write(dataString)
	return err

}

func (s MakeshiftStorage) AddURLPair(shortenedURL string, fullURL string, uuidStr string) {
	if _, found := s.uuidList[uuidStr]; found {
		return
	}
	s.AddUrlPairInMemory(shortenedURL, fullURL, uuidStr)
	writeIntoFile(configuration.ReadFlags().FileStoragePath, DataStorageRecord{UUID: uuidStr,
		ShortUrl: strings.TrimPrefix(
			strings.TrimPrefix(shortenedURL, configuration.ReadFlags().ShortenerBaseURL), "/"), OriginalUrl: fullURL})
}

func (s MakeshiftStorage) AddUrlPairInMemory(shortenedURL string, fullURL string, uuidStr string) {
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
