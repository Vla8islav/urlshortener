package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/errcustom"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
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

	_, err = instance.connPool.Exec(ctx, "CREATE TABLE IF NOT EXISTS url_mapping (UUID char(36) PRIMARY KEY, ShortURL varchar(2000), OriginalURL varchar(2000), UserID integer, Deleted boolean DEFAULT FALSE)")
	if err != nil {
		panic("Couldn't create postgres table" + err.Error())
	}

	_, err = instance.connPool.Exec(ctx, "CREATE TABLE IF NOT EXISTS users (UserID SERIAL PRIMARY KEY)")
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

func (s PostgresStorage) AddURLPairInMemory(ctx context.Context, shortenedURL string, fullURL string, uuidStr string, userID int) {
	s.AddURLPair(ctx, shortenedURL, fullURL, uuidStr, userID)
}

type urlMappingTableRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      int    `json:"user_id"`
	Deleted     bool   `json:"deleted"`
}

type usersTableRecord struct {
	UserID sql.NullInt32 `json:"userid"`
}

func (s PostgresStorage) GetFullURL(ctx context.Context, shortenedURL string) (string, error) {

	row := s.connPool.QueryRow(ctx, "SELECT uuid, shorturl, originalurl, deleted FROM "+urlMappingTableName+" WHERE shorturl = $1  LIMIT 1", shortenedURL)
	var u urlMappingTableRecord
	err := row.Scan(&u.UUID, &u.ShortURL, &u.OriginalURL, &u.Deleted)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errcustom.ErrURLNotFound
	} else if err == nil {
		if !u.Deleted {
			return u.OriginalURL, nil
		} else {
			return "", errcustom.ErrURLDeleted
		}
	} else {
		panic(err)
	}

}

func (s PostgresStorage) GetShortenedURL(ctx context.Context, fullURL string) (string, int, bool) {
	row := s.connPool.QueryRow(ctx, "SELECT uuid, shorturl, originalurl, userid FROM "+urlMappingTableName+" WHERE originalurl = $1  LIMIT 1", fullURL)
	var u urlMappingTableRecord
	err := row.Scan(&u.UUID, &u.ShortURL, &u.OriginalURL, &u.UserID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", -1, false
	} else if err == nil {
		return u.ShortURL, u.UserID, true
	} else {
		panic(err)
	}
}

func (s PostgresStorage) Ping(ctx context.Context) error {
	_, err := s.connPool.Exec(ctx, "select * from urlshortener.public.url_mapping limit 1")
	return err
}

func (s PostgresStorage) GetAllURLRecordsByUser(ctx context.Context, userID int) ([]URLPair, error) {
	rows, err := s.connPool.Query(ctx, "SELECT shorturl, originalurl FROM "+urlMappingTableName+" WHERE userid = $1", userID)
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

func (s PostgresStorage) GetNewUserID(ctx context.Context) (int, error) {
	row := s.connPool.QueryRow(ctx, "insert into users(userid) values(default) returning userid")

	var u usersTableRecord
	err := row.Scan(&u.UserID)

	if errors.Is(err, pgx.ErrNoRows) {
		return -1, fmt.Errorf("couldn't get user id")
	} else if err == nil {
		return int(u.UserID.Int32), nil
	} else {
		panic(err)
	}

}

func (s PostgresStorage) DeleteURL(ctx context.Context, shortenedURL string) error {
	url, _ := helpers.ShortKeyToURL(shortenedURL)
	r, err := s.connPool.Exec(ctx, "update url_mapping set deleted = TRUE where shorturl = $1", url)
	if err != nil {
		return err
	}
	if r.RowsAffected() == 0 {
		return errors.New("couldn't find a requested URL " + url)
	}
	return err

}
