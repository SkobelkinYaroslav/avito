package errGroup

import "errors"

var NotFound = errors.New("not found")
var AlreadyExists = errors.New("already exists")
