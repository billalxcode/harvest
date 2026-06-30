package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"harvest/internal/session"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sardanioss/httpcloak"
)

const googleNewsDefaultRPCID = "Fbv4je"
const googleNewsBaseURL = "https://news.google.com"

var patterns = regexp.MustCompile(`"wrb.fr","Fbv4je","\[\\{,2}"garturlres\\{,2}",\\{,2}"(https?://.*?)\\{,2}"`)
var defaultHeaders = http.Header{"Host": {"news.google.com"}, "Origin": {googleNewsBaseURL}, "Accept-Language": []string{"en-US"}}

type Adapter struct {
	client *httpcloak.Session
}

func NewGoogleAdapter(sessionManager *session.Manager) *Adapter {
	client, err := sessionManager.NewClient()
	if err != nil {
		panic(err)
	}

	return &Adapter{
		client: client,
	}
}

func (a *Adapter) Resolve(id string) (string, error) {
	response, err := a.GetResponse(id)
	if err != nil {
		return "", err
	}
	payload, err := a.BuildRSSPayload(response)
	if err != nil {
		return "", err
	}

	originURL, err := a.TryExtractSource(payload, id)
	if err != nil {
		return "", err
	}

	return originURL, nil
}

func (a *Adapter) GetResponse(id string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	articleURL := a.BuildArticleURL(id)
	response, err := a.client.Get(ctx, articleURL)
	if err != nil {
		panic(err)
	}

	return response.Text()
}

func (a *Adapter) TryExtractSource(payload string, id string) (string, error) {
	batchURL := a.BuildBatchExecuteURL(payload, id)

	headers := a.BuildHeaders()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := httpcloak.Request{
		Method:  "POST",
		URL:     batchURL,
		Body:    strings.NewReader(fmt.Sprintf("f.req=%s", payload)),
		Headers: headers,
	}
	response, err := a.client.Do(ctx, &request)
	if err != nil {
		return "", err
	}

	responseText, err := response.Text()
	if err != nil {
		return "", err
	}

	match := patterns.FindStringSubmatch(responseText)
	if len(match) < 2 {
		return "", nil
	}
	return match[1], nil
}

func (a *Adapter) BuildHeaders() http.Header {
	headers := defaultHeaders.Clone()
	headers.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	headers.Set("Accept-Encoding", "gzip, deflate, br")

	return headers
}

func (a *Adapter) BuildArticleURL(id string) string {
	articleURL := url.URL{
		Scheme: "https",
		Host:   "news.google.com",
		Path:   "/rss/articles",
	}
	articleURL = *articleURL.JoinPath(id)

	return articleURL.String()
}

func (a *Adapter) BuildBatchExecuteURL(payload string, id string) string {
	batchURL := url.URL{
		Scheme: "https",
		Host:   "news.google.com",
		Path:   "/_/DotsSplashUi/data/batchexecute"}
	urlQuery := batchURL.Query()
	urlQuery.Set("rpcids", googleNewsDefaultRPCID)
	urlQuery.Set("source-path", fmt.Sprintf("/rss/articles/%s", id))

	batchURL.RawQuery = urlQuery.Encode()

	return batchURL.String()
}

func (a *Adapter) BuildRSSPayload(rssResponse string) (string, error) {
	document, err := goquery.NewDocumentFromReader(
		strings.NewReader(rssResponse))
	if err != nil {
		return "", err
	}
	dataP, exists := document.Find("c-wiz").First().Attr("data-p")
	if !exists || dataP == "" {
		return "", errors.New("missing 'c-wiz' and 'data-p' attribute")
	}
	dataPReplaced := strings.Replace(dataP, "%.@.", `["garturlreq",`, 1)
	var dataPJSON []any
	if err := json.Unmarshal([]byte(dataPReplaced), &dataPJSON); err != nil {
		return "", fmt.Errorf("unmarshal data-p payload: %w", err)
	}

	patchDataPTimestamp(dataPJSON)
	dataPModified := truncateDataP(dataPJSON)

	dataPBytes, err := json.Marshal(dataPModified)
	if err != nil {
		return "", fmt.Errorf("marshal modified data-p payload: %w", err)
	}

	dataPString := strings.NewReplacer("false", "0", "true", "1").Replace(string(dataPBytes))

	rssPayloadRaw := []any{
		[]any{
			[]any{googleNewsDefaultRPCID, dataPString, nil, "generic"},
		},
	}

	rssPayloadBytes, err := json.Marshal(rssPayloadRaw)
	if err != nil {
		return "", fmt.Errorf("marshal rss payload: %w", err)
	}

	return string(rssPayloadBytes), nil
}

func patchDataPTimestamp(dataPJSON []any) {
	if len(dataPJSON) <= 1 {
		return
	}

	outer, ok := dataPJSON[1].([]any)
	if !ok || len(outer) == 0 {
		return
	}

	inner, ok := outer[0].([]any)
	if !ok || len(inner) <= 9 {
		return
	}

	inner[9] = 420
}

func truncateDataP(dataPJSON []any) []any {
	if len(dataPJSON) < 5 {
		return dataPJSON
	}

	truncated := make([]any, 0, 5)
	truncated = append(truncated, dataPJSON[:3]...)
	truncated = append(truncated, dataPJSON[len(dataPJSON)-2:]...)

	return truncated
}
