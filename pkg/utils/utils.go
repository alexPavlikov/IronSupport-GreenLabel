package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
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

func ReadCookies(r *http.Request) []string {
	var c *http.Cookie
	for _, c = range r.Cookies() {
		if c.Value != "" {
			arr := strings.Split(c.Value, " ")
			return arr
		} else {
			return nil
		}
	}
	return nil
}

func CreateMd5Hash(text string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, text)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
