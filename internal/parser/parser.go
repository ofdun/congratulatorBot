package parser

import (
	"CongratulatorBot/internal/model"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
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
	defer func() {
		if err = response.Body.Close(); err != nil {
			return
		}
	}()

	html, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	p.HTML = string(html)

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
	var holidayTitle, holidayDay, holidayHref string

	document, err := goquery.NewDocumentFromReader(strings.NewReader(p.HTML))
	if err != nil {
		return holidays, err
	}

	document.Find(".album-info").Each(func(i int, s *goquery.Selection) {

		s.Find(".name").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				holidayTitle = s.Text()
				holidayHref, _ = s.Attr("href")
			} else if i == 2 {
				UnparsedDate := s.Text()
				holidayDay = ParseDateFromString(UnparsedDate)
			}
		})
		holiday := model.NewHoliday(
			holidayTitle,
			holidayDay,
			holidayHref,
		)
		holidays = append(holidays, *holiday)
	})

	return holidays, err
}

func (p *Parser) GetPostcardsPages(holiday model.Holiday) ([]model.Postcard, error) {
	var postcards []model.Postcard
	var postcardHref string

	response, err := http.Get(holiday.Href)
	if err != nil {
		return postcards, err
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			return
		}
	}()

	html, err := io.ReadAll(response.Body)
	if err != nil {
		return postcards, err
	}

	HTMLPage := string(html)
	document, err := goquery.NewDocumentFromReader(strings.NewReader(HTMLPage))

	if err != nil {
		return postcards, err
	}

	document.Find(".card-image").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				postcardHref = href
				postcard := model.NewPostcard(
					holiday.Name,
					postcardHref,
				)
				postcards = append(postcards, *postcard)
			}
		})
	})

	return postcards, err
}

func youTubeLink(url string) bool {
	return strings.Contains(url, "youtube")
}

func (p *Parser) GetPostcardHref(postcard *model.Postcard) error {
	response, err := http.Get(postcard.Page)
	defer func() {
		if err = response.Body.Close(); err != nil {
			return
		}
	}()

	if err != nil {
		return err
	}

	html, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	document, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))

	document.Find(".cardContent").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			if youTubeLink(src) {
				postcard.YouTube = true
			}
			postcard.Href = src
		} else {
			s.Find("source").Each(func(i int, s *goquery.Selection) {
				src, exists = s.Attr("src")
				if exists {
					postcard.Href = src
				}
			})
		}
	})

	if err != nil {
		return err
	}

	return nil
}
