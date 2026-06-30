package main

import (
	"fmt"
	"harvest/internal/core"
	"harvest/internal/engine"
	"harvest/internal/engine/instances"
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

	result, err := engineManager.Search("google", "majalengka")
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	_ = result
}
