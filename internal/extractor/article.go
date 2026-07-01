package extractor

import (
	"bytes"
	"context"
	"fmt"
	"harvest/internal/session"
	"harvest/internal/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"codeberg.org/readeck/go-readability/v2"
	"github.com/markusmobius/go-trafilatura"
	"github.com/sardanioss/httpcloak"
)

type ArticleExtractor struct {
	client *httpcloak.Session
}

type TrafilaturaContent struct {
	Title       string
	Author      string
	URL         string
	Hostname    string
	SiteName    string
	Date        string
	Categories  []string
	Tags        []string
	Fingerprint string
	Text        string
	Language    string
	Image       string
	PageType    string
}

type ReadabilityContent struct {
	Title        string
	Author       string
	SiteName     string
	Language     string
	Image        string
	Body         string
	Text         string
	ModifiedDate string
}

func extractWithTrafilatura(raw []byte, pageURL *url.URL) (TrafilaturaContent, error) {
	result, err := trafilatura.Extract(bytes.NewReader(raw), trafilatura.Options{
		OriginalURL:    pageURL,
		EnableFallback: true,
		IncludeImages:  true})
	if err != nil {
		return TrafilaturaContent{}, fmt.Errorf("failed to extract trafilatura: %w", err)
	}

	metadata := result.Metadata
	return TrafilaturaContent{Title: metadata.Title, Author: metadata.Author,
		URL:         metadata.URL,
		Hostname:    metadata.Hostname,
		SiteName:    metadata.Sitename,
		Date:        utils.FormatTime(metadata.Date),
		Categories:  metadata.Categories,
		Tags:        metadata.Tags,
		Fingerprint: metadata.Fingerprint,
		Text:        result.ContentText,
		Language:    metadata.Language,
		Image:       metadata.Image,
		PageType:    metadata.PageType}, nil
}

func extractWithReadability(raw []byte, pageURL *url.URL) (ReadabilityContent, error) {
	article, err := readability.FromReader(bytes.NewReader(raw), pageURL)
	if err != nil {
		return ReadabilityContent{}, fmt.Errorf("failed to read article from reader: %w", err)
	}

	var htmlBody strings.Builder
	if err := article.RenderHTML(&htmlBody); err != nil {
		return ReadabilityContent{}, fmt.Errorf("failed to render html: %w", err)
	}

	var plainText strings.Builder
	if err := article.RenderText(&plainText); err != nil {
		return ReadabilityContent{}, fmt.Errorf("failed to render plain text: %w", err)
	}

	modifiedTime, _ := article.ModifiedTime()

	return ReadabilityContent{
		Title:        article.Title(),
		Author:       article.Byline(),
		SiteName:     article.SiteName(),
		Language:     article.Language(),
		Image:        article.ImageURL(),
		Body:         htmlBody.String(),
		Text:         plainText.String(),
		ModifiedDate: utils.FormatTime(modifiedTime),
	}, nil
}

func extract(body io.Reader, pageURL *url.URL) (ArticleData, error) {
	raw, err := io.ReadAll(body)
	if err != nil {
		return ArticleData{}, fmt.Errorf("failed to read body: %w", err)
	}

	trafilaturaContent, err := extractWithTrafilatura(raw, pageURL)
	if err != nil {
		return ArticleData{}, fmt.Errorf("failed to extract trafilatura: %w", err)
	}

	readabilityContent, err := extractWithReadability(raw, pageURL)
	if err != nil {
		return ArticleData{}, fmt.Errorf("failed to extract readability: %w", err)
	}

	return ArticleData{
		Title:       utils.CoalesceString(trafilaturaContent.Title, readabilityContent.Title),
		Author:      utils.CoalesceString(trafilaturaContent.Author, readabilityContent.Author),
		URL:         trafilaturaContent.URL,
		Hostname:    trafilaturaContent.Hostname,
		SiteName:    utils.CoalesceString(trafilaturaContent.SiteName, readabilityContent.SiteName),
		Date:        trafilaturaContent.Date,
		Categories:  trafilaturaContent.Categories,
		Tags:        trafilaturaContent.Tags,
		Fingerprint: trafilaturaContent.Fingerprint,
		Body:        readabilityContent.Body,
		Text:        utils.CoalesceString(trafilaturaContent.Text, readabilityContent.Text),
		Language:    utils.CoalesceString(trafilaturaContent.Language, readabilityContent.Language),
		Image:       utils.CoalesceString(trafilaturaContent.Image, readabilityContent.Image),
		PageType:    trafilaturaContent.PageType,
		Metadata: ArticleMetadata{
			PublishedDate: trafilaturaContent.Date,
			ModifiedDate:  readabilityContent.ModifiedDate,
			Source:        trafilaturaContent.Hostname,
			Author:        utils.CoalesceString(trafilaturaContent.Author, readabilityContent.Author),
		},
	}, nil
}

func NewArticleExtractor(sessionManager *session.Manager) (*ArticleExtractor, error) {
	client, err := sessionManager.NewClient()
	if err != nil {
		return nil, err
	}

	return &ArticleExtractor{
		client: client,
	}, nil
}

func (ae *ArticleExtractor) Extract(originURL string) (ArticleData, error) {
	parsedURL, err := url.Parse(originURL)
	if err != nil {
		return ArticleData{}, fmt.Errorf("failed to parse URL: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := ae.client.Do(ctx, &httpcloak.Request{
		Method: "GET",
		URL:    parsedURL.String(),
		Headers: http.Header{
			"Accept-Language": {"id,id-ID;q=0.9,en;q=0.5"},
			"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		},
	})
	if err != nil {
		return ArticleData{}, fmt.Errorf("failed to request: %w", err)
	}

	if response.StatusCode >= 400 {
		return ArticleData{}, fmt.Errorf("failed to request, status code %d", response.StatusCode)
	}

	data, err := extract(response.Body, parsedURL)
	if err != nil {
		return ArticleData{}, fmt.Errorf("failed to extract article: %w", err)
	}

	return data, nil
}
