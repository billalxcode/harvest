package engine

type Article struct {
	title       string `validate:"required,min=3,max=100"`
	link        string `validate:"required,min=1"`
	guid        string `validate:"required,min=1"`
	description string `validate:"required,min=1"`
	pubDate     string `validate:"required,min=1"`
	source      string `validate:"required,min=1"`
	originURL   string `validate:"required,min=1"`
}

type SearchResult struct {
	query    string    `validate:"required,min=1"`
	engine   string    `validate:"required,min=1"`
	articles []Article `validate:"optional"`
	total    int       `validate:"required,min=1"`
}

type BaseEngineInstance interface {
	Search(query string) (SearchResult, error)
}
