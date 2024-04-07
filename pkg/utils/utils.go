package utils

import (
	"bufio"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	gomail "gopkg.in/mail.v2"

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

type UtilsProduct struct {
	Id                string
	Name              string
	FullName          string
	Waight            int //gramm
	Category          string
	UnitOfMeasurement string
	Remains           int
	Price             int
}

func WriteProductToExcelFile(products []UtilsProduct) error {
	file, err := excelize.OpenFile("./yandex_disk/DBSynchronization")
	if err != nil {
		return err
	}

	for i := 1; i <= len(products)-1; i++ {
		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), products[i-1].Id)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", i), products[i-1].Name)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%d", i), products[i-1].FullName)
		file.SetCellValue("Sheet1", fmt.Sprintf("D%d", i), products[i-1].Remains)
		file.SetCellValue("Sheet1", fmt.Sprintf("E%d", i), products[i-1].UnitOfMeasurement)
		file.SetCellValue("Sheet1", fmt.Sprintf("F%d", i), products[i-1].Price)
		file.SetCellValue("Sheet1", fmt.Sprintf("G%d", i), products[i-1].Waight)
	}

	err = file.Save()
	if err != nil {
		return err
	}

	return nil
}

func SendMessage(email string, message string, header string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "a.pavlikov2002@gmail.com")
	m.SetHeader("To", email)

	m.SetHeader("Subject", header)

	m.SetBody("text/plain", message)
	d := gomail.NewDialer("smtp.gmail.com", 587, "a.pavlikov2002@gmail.com", "isei dkte iiwl wior")

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
