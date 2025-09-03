package helper

import (
	"time"
)

func PtrTime(t time.Time) *time.Time {
	return &t
}