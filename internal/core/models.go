package core

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ResultItem struct {
	Title       string `validate:"required,min=3"`
	Link        string `validate:"required,url"`
	GUID        string `validate:"required"`
	Description string `validate:"required"`
	PubDate     string `validate:"required"`
	Source      string `validate:"required"`
	OriginURL   string `validate:"required,url"`
}

type SearchResult struct {
	Query  string       `validate:"required"`
	Engine string       `validate:"required"`
	Items  []ResultItem `validate:"dive"`
	Total  int          `validate:"min=0"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func NewSearchItem(
	title string,
	link string,
	guid string,
	description string,
	pubDate string,
	source string,
	originURL string) (*ResultItem, error) {
	item := ResultItem{
		Title:       title,
		Link:        link,
		GUID:        guid,
		Description: description,
		PubDate:     pubDate,
		Source:      source,
		OriginURL:   originURL,
	}
	err := validate.Struct(item)
	if err != nil {
		return nil, fmt.Errorf("invalid item: %s", err.Error())
	}
	return &item, nil
}

func NewSearchResult(query string, engine string) (*SearchResult, error) {
	result := &SearchResult{Query: query, Engine: engine, Items: make([]ResultItem, 0), Total: 0}
	err := validate.Var(result.Query, "required")
	if err != nil {
		return nil, fmt.Errorf("invalid query: %s", err.Error())
	}

	err = validate.Var(result.Engine, "required")
	if err != nil {
		return nil, fmt.Errorf("invalid engine: %s", err.Error())
	}

	return result, nil
}

func (r *SearchResult) AppendItem(item ResultItem) error {
	err := validate.Struct(item)
	if err != nil {
		return fmt.Errorf("invalid item: %s", err.Error())
	}

	r.Items = append(r.Items, item)
	r.Total = len(r.Items)

	return nil
}
