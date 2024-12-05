package service

import "errors"

var (
	ErrHasNotToken = errors.New("user has not token")
	ErrInternal = errors.New("error from repo")
	ErrConcurency = errors.New("smth was, but now no or change")
)
