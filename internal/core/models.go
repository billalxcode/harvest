package core

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Article struct {
	Title       string `validate:"required,min=3"`
	Link        string `validate:"required,url"`
	GUID        string `validate:"required"`
	Description string `validate:"required"`
	PubDate     string `validate:"required"`
	Source      string `validate:"required"`
	OriginURL   string `validate:"required,url"`
}

type SearchResult struct {
	Query    string    `validate:"required"`
	Engine   string    `validate:"required"`
	Articles []Article `validate:"dive"`
	Total    int       `validate:"min=0"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func NewArticle(
	title string,
	link string,
	guid string,
	description string,
	pubDate string,
	source string,
	originURL string) (Article, error) {
	article := Article{
		Title:       title,
		Link:        link,
		GUID:        guid,
		Description: description,
		PubDate:     pubDate,
		Source:      source,
		OriginURL:   originURL,
	}
	err := validate.Struct(article)
	if err != nil {
		return Article{}, fmt.Errorf("invalid article: %s", err.Error())
	}
	return article, nil
}

func NewSearchResult(query string, engine string) (*SearchResult, error) {
	result := &SearchResult{Query: query, Engine: engine, Articles: make([]Article, 0), Total: 0}
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

func (r *SearchResult) AppendArticle(article Article) error {
	err := validate.Struct(article)
	if err != nil {
		return fmt.Errorf("invalid article: %s", err.Error())
	}

	r.Articles = append(r.Articles, article)
	r.Total = len(r.Articles)

	return nil
}
