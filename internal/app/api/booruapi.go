package api

import (
	"context"
	"errors"
	"github.com/Vla8islav/urlshortener/internal/app/errcustom"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
)

type GoBooruServiceMethods interface {
	Register(ctx context.Context, username string, password string) (string, error)
}

type GoBooruService struct {
	Storage storage.Storage
}

func NewGoBooruService(ctx context.Context, s storage.Storage) (GoBooruService, error) {
	return GoBooruService{Storage: s}, nil
}

func (u GoBooruService) Register(ctx context.Context, username string, password string) (int, error) {
	var err error
	var userID int

	if username == "" {
		return -1, errors.New("incorrect username")
	} else if userID, err = u.Storage.GetUserByUsername(ctx, username); !errors.Is(err, errcustom.ErrUserNotFound) {
		return userID, errcustom.ErrUserAlreadyExists
	}
	if password == "" {
		return -1, errors.New("incorrect password")
	}

	userID, err = u.Storage.CreateUser(ctx, username, password)
	return userID, err
}
