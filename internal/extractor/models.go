package extractor

type ArticleMetadata struct {
	PublishedDate string `json:"published_date"`
	ModifiedDate  string `json:"modified_date"`
	Source        string `json:"source"`
	RawPublished  string `json:"raw_published"`
	RawModified   string `json:"raw_modified"`
	Confidence    string `json:"confidence"`
	Author        string `json:"author"`
	RawAuthor     string `json:"raw_author"`
}

type ArticleData struct {
	Title       string          `json:"title"`
	Author      string          `json:"author"`
	URL         string          `json:"url"`
	Hostname    string          `json:"hostname"`
	SiteName    string          `json:"site_name"`
	Date        string          `json:"date"`
	Categories  []string        `json:"categories"`
	Tags        []string        `json:"tags"`
	Fingerprint string          `json:"fingerprint"`
	Body        string          `json:"body"`
	Text        string          `json:"text"`
	Language    string          `json:"language"`
	Image       string          `json:"image"`
	PageType    string          `json:"page_type"`
	Metadata    ArticleMetadata `json:"metadata"`
}
