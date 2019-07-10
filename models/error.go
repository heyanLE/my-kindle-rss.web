package models

import "errors"

var (
	IncorrectUserOPassErr = errors.New("Incorrect username or password ")
	EmailExistErr         = errors.New("Exist Email ")
	FeedNotFound          = errors.New("Feed Not Found ")
	UserNotFound          = errors.New("User Not Found ")
)
