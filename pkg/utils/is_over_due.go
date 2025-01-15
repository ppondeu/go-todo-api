package utils

import "time"

func IsOverDue(dueDate time.Time) bool {
	return time.Now().After(dueDate)
}
