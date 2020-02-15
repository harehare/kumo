package kumo

import "errors"

var ErrNotAllowDomain = errors.New("Not allow domain")
var ErrRobotsTxtBlocked = errors.New("Robots.txt blocked")
var ErrInvalidParameter = errors.New("Invalid Parameter")
