package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product"
	"github.com/xuri/excelize/v2"
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

func ReadEventFile() (lines []string, err error) {
	path := "./logs/events.log"
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		} else {
			lines = append(lines, scanner.Text())
		}
	}

	var arr []string

	for i := len(lines) - 1; i > 0; i-- {
		arr = append(arr, lines[i])
	}

	return arr, scanner.Err()
}

func WriteProductToExcelFile(products []product.Product) error {
	file, err := excelize.OpenFile("./yandex_disk/DBSynchronization")
	if err != nil {
		return err
	}

	for i := 1; i <= len(products)-1; i++ {
		file.SetCellValue("Sheet1", "A"+string(i), "")
		file.SetCellValue("Sheet1", "B"+string(i), "")
		file.SetCellValue("Sheet1", "C"+string(i), "")
		file.SetCellValue("Sheet1", "D"+string(i), "")
		file.SetCellValue("Sheet1", "E"+string(i), "")
	}

	return nil
}
