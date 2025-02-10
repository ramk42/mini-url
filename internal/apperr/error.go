package apperr

import "errors"

var (
	ErrURLNotFound = errors.New("URL not found")

	ErrURLUpdateClicksFailed = errors.New("URL update clicks failed")

	ErrURLFetchFailed = errors.New("URL fetch failed")

	ErrURLInsertFailed = errors.New("URL insert failed")

	ErrURLConflict = errors.New("URL insert conflict")
)
