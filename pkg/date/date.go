package date

import (
	"time"
)

func PtrTime(t time.Time) *time.Time {
	return &t
}