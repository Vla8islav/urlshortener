package errlist

import "errors"

var ErrURLNotFound = errors.New("couldn't find a requested URLRepo")
