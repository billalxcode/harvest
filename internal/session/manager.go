package session

import (
	"harvest/internal/utils"
	"time"

	"github.com/sardanioss/httpcloak"
)

type Proxy struct {
	http string `validate:"required,url"`
}

type Manager struct {
	client  *httpcloak.Session
	proxies *utils.Cycler[Proxy]
}

func NewManager(proxies []Proxy) (*Manager, error) {
	pool := utils.NewCycler(proxies)

	return &Manager{
		client:  nil,
		proxies: pool,
	}, nil
}

func (sm *Manager) NewClient() (*httpcloak.Session, error) {
	client := httpcloak.NewSession("chrome-latest", httpcloak.WithSessionTimeout(30*time.Second))

	if sm.proxies != nil && sm.proxies.Length() > 0 {
		proxy := sm.GetNextProxy()
		client.SetProxy(proxy.http)
	}

	return client, nil
}

func (sm *Manager) GetNextProxy() Proxy {
	proxy, err := sm.proxies.Next()
	if err != nil {
		panic(err)
	}
	return proxy
}
