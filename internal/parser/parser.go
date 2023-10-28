package parser

import (
	"io"
	"net/http"
)

type Parser struct {
	URL  string
	HTML string
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
