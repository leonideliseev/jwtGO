package service

import "errors"

var (
	HasNotToken = errors.New("user has not token")
	ErrConcurency = errors.New("smth was, but now no or change")
)
