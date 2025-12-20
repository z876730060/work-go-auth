package feign

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	FeignClients = NewClientManager()
	httpClient   = http.DefaultClient
)

type DiscoverInstance struct {
	Scheme      string `json:"scheme"`
	Ip          string `json:"ip"`
	Port        uint64 `json:"port"`
	ServiceName string `json:"serviceName"`
}

type ClientManager struct {
	clientMap      map[string]DiscoverInstance
	registerClient map[string]bool
	mux            sync.Mutex
}

func (c *ClientManager) NeedDiscover(serviceName string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.registerClient[serviceName]
	return ok
}

func (c *ClientManager) Clear() {
	c.mux.Lock()
	defer c.mux.Unlock()
	clear(c.clientMap)
}

func (c *ClientManager) Register(instance DiscoverInstance) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.clientMap[instance.ServiceName] = instance
}

func (c *ClientManager) register(s string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.registerClient[s] = true
}

func (c *ClientManager) get(s string) (discoverInstance DiscoverInstance, err error) {
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
	register(string)
	get(string) (http.Client, error)
	Register(DiscoverInstance)
	NeedDiscover(string) bool
}

type Discover interface {
	Discover(interface{})
}
