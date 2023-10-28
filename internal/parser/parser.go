package parser

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"telegramBot/internal/model"
	"unicode"
)

type Parser struct {
	URL  string
	HTML string
}

func NewParser(url string) *Parser {
	return &Parser{URL: url}
}

func (p *Parser) GetHTML() error {
	response, err := http.Get(p.URL)
	if err != nil {
		return err
	}

	html, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	p.HTML = string(html)

	err = response.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func ParseDateFromString(date string) string {
	newDate := ""
	for _, char := range date {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			newDate += string(char)
		} else if char == ' ' {
			break
		}
	}
	return newDate
}

func (p *Parser) GetHolidays() ([]model.Holiday, error) {
	var holidays []model.Holiday
	var holidayTitle, holidayDay string

	document, err := goquery.NewDocumentFromReader(strings.NewReader(p.HTML))
	if err != nil {
		return holidays, err
	}

	document.Find(".album-info").Each(func(i int, s *goquery.Selection) {

		s.Find(".name").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				holidayTitle = s.Text()
			} else if i == 2 {
				UnparsedDate := s.Text()
				holidayDay = ParseDateFromString(UnparsedDate)
			}
		})
		holiday := model.NewHoliday(
			holidayTitle,
			holidayDay,
		)
		holidays = append(holidays, *holiday)
	})

	return holidays, err
}
