package model

import "errors"

var (
	ErrNotFound         = errors.New("subscription not found")
	ErrInvalidUUID      = errors.New("invalid user_id (not a valid UUID")
	ErrInvalidDateRange = errors.New("end_date must be after start_date")
)
