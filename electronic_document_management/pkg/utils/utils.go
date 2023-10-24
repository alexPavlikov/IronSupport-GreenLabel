package utils

import (
	"strings"
	"time"
)

func DoWithTries(fn func() error, attempts int, duration time.Duration) (err error) {
	for attempts > 0 {
		err = fn()
		if err != nil {
			time.Sleep(duration)
			attempts--
			continue
		}
		return nil
	}
	return err
}

func FormatQuery(query string) string {
	return strings.ReplaceAll(strings.ReplaceAll(query, "\t", " "), "\n", " ")
}
