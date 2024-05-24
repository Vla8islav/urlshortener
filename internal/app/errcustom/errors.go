package errcustom

import (
	"errors"
)

var ErrURLNotFound = errors.New("couldn't find a requested URL")
var ErrURLDeleted = errors.New("a requested URL is deleted")
var ErrUserNotFound = errors.New("couldn't find user with that ID")
var ErrUserAlreadyExists = errors.New("user with that ID already exists")
