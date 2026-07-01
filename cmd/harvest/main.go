package main

import (
	"fmt"
	"harvest/internal/core"
	"harvest/internal/engine"
	"harvest/internal/engine/instances"
	"harvest/internal/extractor"
	"harvest/internal/session"
)

func main() {
	config := core.Config{}
	sessionManager, err := session.NewManager(nil)
	if err != nil {
		panic(err)
	}

	googleEngine := instances.NewGoogleEngine(
		config,
		sessionManager,
	)
	engineManager := engine.NewManager()
	engineManager.RegisterEngine("google", googleEngine)

	results, err := engineManager.Search("google", "majalengka")
	if err != nil {
		panic(err)
	}

	extractor_, err := extractor.NewArticleExtractor(sessionManager)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Extract data from %s", results.Items[3].OriginURL))
	content, err := extractor_.Extract(results.Items[3].OriginURL)
	if err != nil {
		panic(err)
	}

	fmt.Println(content)

	_ = results
}
