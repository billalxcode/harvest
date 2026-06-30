package instances

import (
	"context"
	"fmt"
	"harvest/internal/core"
	"harvest/internal/engine"
	"harvest/internal/parser"
	"harvest/internal/session"
	"net/url"
	"time"
)

type GoogleEngine struct {
	config         core.Config
	sessionManager *session.Manager
}

func NewGoogleEngine(config core.Config, sessionManager *session.Manager) *GoogleEngine {
	return &GoogleEngine{
		config:         config,
		sessionManager: sessionManager,
	}
}

func (ge *GoogleEngine) Search(query string) (engine.SearchResult, error) {
	rssURL := ge.InitializeURL(query)

	client, err := ge.sessionManager.NewClient()
	if err != nil {
		panic(err)
	}

	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := client.Get(ctx, rssURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("status", response.StatusCode)
	body, err := response.Bytes()
	if err != nil {
		return engine.SearchResult{}, err
	}

	content, err := parser.Decode(body)
	if err != nil {
		return engine.SearchResult{}, err
	}

	fmt.Println("content", content)

	return engine.SearchResult{}, nil
}

func (ge *GoogleEngine) InitializeURL(query string) string {
	rssURL := url.URL{
		Scheme: "https",
		Host:   "news.google.com",
		Path:   "/rss/search",
	}
	urlQuery := rssURL.Query()
	urlQuery.Set("q", query)
	urlQuery.Set("hl", ge.config.Language)
	urlQuery.Set("gl", ge.config.Country)
	urlQuery.Set("ceid", fmt.Sprintf("%s:%s", ge.config.Language, ge.config.Country))

	rssURL.RawQuery = urlQuery.Encode()

	return rssURL.String()
}
