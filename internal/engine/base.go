package engine

import "harvest/internal/core"

type BaseEngineInstance interface {
	Search(query string) (*core.SearchResult, error)
}
