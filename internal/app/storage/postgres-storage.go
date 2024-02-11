package storage

import (
	"context"
	"errors"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/jackc/pgx/v5"
)

const urlMappingTableName = "url_mapping"

type PostgresStorage struct {
}

func NewPostgresStorage() (Storage, error) {
	instance := new(PostgresStorage)
	return instance, nil
}

func (s PostgresStorage) AddURLPair(shortenedURL string, fullURL string, uuidStr string) {
	ctx, conn := getPostgresConnection()
	defer conn.Close(ctx)

	_, err := conn.Exec(ctx, "INSERT INTO "+urlMappingTableName+"(UUID, ShortURL, OriginalURL) values ($1, $2, $3)", uuidStr, shortenedURL, fullURL)
	if err != nil {
		panic("Couldn't insert data into" + urlMappingTableName + " postgres table")
	}

}

func getPostgresConnection() (context.Context, *pgx.Conn) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, configuration.ReadFlags().DBConnectionString)
	if err != nil {
		panic("Couldn't create connection to the postgres DB")
	}
	return ctx, conn
}

func (s PostgresStorage) AddURLPairInMemory(shortenedURL string, fullURL string, uuidStr string) {
	s.AddURLPair(shortenedURL, fullURL, uuidStr)
}

type urlMappingTableRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s PostgresStorage) GetFullURL(shortenedURL string) (string, bool) {
	ctx, conn := getPostgresConnection()
	defer conn.Close(ctx)

	row := conn.QueryRow(ctx, "SELECT uuid, shorturl, originalurl FROM "+urlMappingTableName+" WHERE shorturl = $1  LIMIT 1", shortenedURL)
	var u urlMappingTableRecord
	err := row.Scan(&u.UUID, &u.ShortURL, &u.OriginalURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false
	} else if err == nil {
		return u.OriginalURL, true
	} else {
		panic(err)
	}

}

func (s PostgresStorage) GetShortenedURL(fullURL string) (string, bool) {
	ctx, conn := getPostgresConnection()
	defer conn.Close(ctx)

	row := conn.QueryRow(ctx, "SELECT uuid, shorturl, originalurl FROM "+urlMappingTableName+" WHERE originalurl = $1  LIMIT 1", fullURL)
	var u urlMappingTableRecord
	err := row.Scan(&u.UUID, &u.ShortURL, &u.OriginalURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false
	} else if err == nil {
		return u.ShortURL, true
	} else {
		panic(err)
	}
}
