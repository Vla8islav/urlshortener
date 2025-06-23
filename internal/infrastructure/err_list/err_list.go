package err_list

import "errors"

var ErrURLNotFound = errors.New("couldn't find a requested URLRepo")
