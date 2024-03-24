package errcustom

import "errors"

var ErrURLNotFound = errors.New("couldn't find a requested URL")
var ErrURLDeleted = errors.New("a requested URL is deleted")
