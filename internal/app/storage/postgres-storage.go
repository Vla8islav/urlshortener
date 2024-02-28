package storage

import (
	"context"
	"errors"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const urlMappingTableName = "url_mapping"

type PostgresStorage struct {
	connPool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context) (Storage, error) {
	instance := new(PostgresStorage)
	var err error
	instance.connPool, err = pgxpool.New(ctx, configuration.ReadFlags().DBConnectionString)
	if err != nil {
		panic("Couldn't connect to the postgres server" + err.Error())
	}

	_, err = instance.connPool.Exec(ctx, "CREATE TABLE IF NOT EXISTS url_mapping (UUID char(36) PRIMARY KEY, ShortURL varchar(2000), OriginalURL varchar(2000), UserID integer)")
	if err != nil {
		panic("Couldn't create postgres table" + err.Error())
	}

	return instance, nil
}

func (s PostgresStorage) Close() {
	s.connPool.Close()
}

func (s PostgresStorage) AddURLPair(ctx context.Context, shortenedURL string, fullURL string, uuidStr string, userID int) {
	_, err := s.connPool.Exec(ctx, "INSERT INTO "+urlMappingTableName+"(UUID, ShortURL, OriginalURL, userid) values ($1, $2, $3, $4)", uuidStr, shortenedURL, fullURL, userID)
	if err != nil {
		panic("Couldn't insert data into" + urlMappingTableName + " postgres table")
	}

}

func getPostgresConnection(ctx context.Context) (context.Context, *pgx.Conn) {
	conn, err := pgx.Connect(ctx, configuration.ReadFlags().DBConnectionString)
	if err != nil {
		panic("Couldn't create connection to the postgres DB")
	}
	return ctx, conn
}

func (s PostgresStorage) AddURLPairInMemory(ctx context.Context, shortenedURL string, fullURL string, uuidStr string, userID int) {
	s.AddURLPair(ctx, shortenedURL, fullURL, uuidStr, userID)
}

type urlMappingTableRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s PostgresStorage) GetFullURL(ctx context.Context, shortenedURL string) (string, bool) {

	row := s.connPool.QueryRow(ctx, "SELECT uuid, shorturl, originalurl FROM "+urlMappingTableName+" WHERE shorturl = $1  LIMIT 1", shortenedURL)
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

func (s PostgresStorage) GetShortenedURL(ctx context.Context, fullURL string) (string, bool) {
	ctx, conn := getPostgresConnection(ctx)
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

func (s PostgresStorage) Ping(ctx context.Context) error {
	_, err := s.connPool.Exec(ctx, "select * from urlshortener.public.url_mapping limit 1")
	return err
}

func (s PostgresStorage) GetAllURLRecordsByUser(ctx context.Context, userId int) ([]URLPair, error) {
	rows, err := s.connPool.Query(ctx, "SELECT shorturl, originalurl FROM "+urlMappingTableName+" WHERE userid = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rowSlice []URLPair
	for rows.Next() {
		var r URLPair
		err = rows.Scan(&r.ShortURL, &r.FullURL)
		if err != nil {
			return nil, err
		}
		rowSlice = append(rowSlice, r)
	}
	return rowSlice, nil
}
