package common

import "errors"

var NotFound = errors.New("NotFound")
var AlreadyExistsErr = errors.New("AlreadyExistsErr")
var Unsupported = errors.New("Unsupported")
var InvalidArgument = errors.New("InvalidArgument")
