package controllers

import (
	"fmt"
	"time"
)

// parseDateString mencoba mem-parse string tanggal dalam beberapa format umum.
// Format yang didukung:
//   - "2006-01-02T15:04:05Z"   (RFC3339 / ISO 8601 UTC)
//   - "2006-01-02T15:04:05+07:00"  (RFC3339 dengan timezone)
//   - "2006-01-02 15:04:05"    (format SQL)
//   - "2006-01-02"             (hanya tanggal)
func parseDateString(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,          // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("format tanggal '%s' tidak dikenali", dateStr)
}
