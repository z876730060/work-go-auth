package feign

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	fClients               = NewClientManager()
	ServiceDiscoverManager = NewDiscoverManager()
	httpClient             = http.DefaultClient
)

type DiscoverInstance struct {
	Scheme      string `json:"scheme"`
	Ip          string `json:"ip"`
	Port        uint64 `json:"port"`
	ServiceName string `json:"serviceName"`
}

type DiscoverManager struct {
	discoverMap map[string][]DiscoverInstance
	mux         sync.Mutex
}

func NewDiscoverManager() *DiscoverManager {
	return &DiscoverManager{
		discoverMap: make(map[string][]DiscoverInstance),
	}
}

func (d *DiscoverManager) Discover(i interface{}) {
	//TODO 先处理DiscoverManager的数据

	//TODO 再去处理FClient的数据
}

type ClientManager struct {
	clientMap      map[string]DiscoverInstance
	registerClient map[string]bool
	mux            sync.Mutex
}

func (c *ClientManager) Register(s string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.registerClient[s] = true
}

func (c *ClientManager) Get(s string) (discoverInstance DiscoverInstance, err error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	var ok bool
	discoverInstance, ok = c.clientMap[s]
	if !ok {
		return discoverInstance, fmt.Errorf("not found %s", s)
	}
	return discoverInstance, nil
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clientMap:      make(map[string]DiscoverInstance),
		registerClient: make(map[string]bool),
	}
}

type FeignClient interface {
	Register(string)
	Get(string) (http.Client, error)
}

type Discover interface {
	Discover(interface{})
}
