package engine

import "fmt"

type Manager struct {
	engines map[string]BaseEngineInstance
}

func NewManager() *Manager {
	return &Manager{
		engines: make(map[string]BaseEngineInstance),
	}
}

func (m *Manager) RegisterEngine(name string, instance BaseEngineInstance) {
	m.engines[name] = instance
}

func (m *Manager) GetEngine(name string) (BaseEngineInstance, error) {
	engine, ok := m.engines[name]
	if !ok {
		return nil, fmt.Errorf("engine '%s' not found", name)
	}
	return engine, nil
}

func (m *Manager) Search(name string, query string) (SearchResult, error) {
	engine, err := m.GetEngine(name)
	if err != nil {
		return SearchResult{}, err
	}
	return engine.Search(query)
}
