package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/errcustom"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	connPool *pgxpool.Pool
}

type UsersTableRecord struct {
	UserID   sql.NullInt32  `json:"userid"`
	Username sql.NullString `json:"username"`
	Password sql.NullString `json:"password"`
	Deleted  sql.NullInt32  `json:"deleted"`
}

func NewPostgresStorage(ctx context.Context) (Storage, error) {
	instance := new(PostgresStorage)
	var err error
	instance.connPool, err = pgxpool.New(ctx, configuration.ReadFlags().DBConnectionString)
	if err != nil {
		panic("Couldn't connect to the postgres server" + err.Error())
	}

	_, err = instance.connPool.Exec(ctx, "CREATE TABLE IF NOT EXISTS users (UserID SERIAL PRIMARY KEY, Username varchar(20), Password varchar(100), Deleted bool DEFAULT FALSE)")
	if err != nil {
		panic("Couldn't create postgres table" + err.Error())
	}

	return instance, nil
}

func (s PostgresStorage) Close() {
	s.connPool.Close()
}

func (s PostgresStorage) Ping(ctx context.Context) error {
	_, err := s.connPool.Exec(ctx, "select * from url_mapping limit 1")
	return err
}

func (s PostgresStorage) CreateUser(ctx context.Context, username string, password string) (int, error) {
	row := s.connPool.QueryRow(ctx, "insert into users(UserID, Username, Password, Deleted) values(default, $1, $2, default) returning userid", username, password)

	var u UsersTableRecord
	err := row.Scan(&u.UserID)

	if errors.Is(err, pgx.ErrNoRows) {
		return -1, fmt.Errorf("couldn't create new user id")
	} else if err == nil {
		return int(u.UserID.Int32), nil
	} else {
		panic(err)
	}
}

func (s PostgresStorage) GetUserByUsername(ctx context.Context, username string) (int, error) {
	row := s.connPool.QueryRow(ctx, "select UserID from users where UserID = $1 ", username)

	var u UsersTableRecord
	err := row.Scan(&u.UserID)

	if errors.Is(err, pgx.ErrNoRows) {
		return -1, errcustom.ErrUserNotFound
	} else if err == nil {
		return int(u.UserID.Int32), nil
	} else {
		panic(err)
	}
}
