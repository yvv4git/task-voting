package service

import "errors"

var (
	ErrStartAtBeforeEndAt = errors.New("start_at must be before end_at")
)
