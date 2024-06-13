package handler

import (
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Client struct {
	Conn   *websocket.Conn
	UUID   string
	Domain string
	Locker sync.Mutex // 锁，同一个 client 发消息必须加锁，不然并发的时候报错
}

type Manager struct {
	Clients map[string]*Client   // uuid:client
	Domains map[string][]*Client // domain:[client]
	Locker  sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		Clients: map[string]*Client{},
		Domains: map[string][]*Client{},
	}
}

// Exists 判断域名是否存在
func (m *Manager) ExistsByDomain(domain string) bool {
	m.Locker.RLock()
	defer m.Locker.RUnlock()
	if _, ok := m.Domains[domain]; ok {
		if len(m.Domains[domain]) != 0 {
			return true
		}
	}
	return false
}

// Exists 判断域名是否存在
func (m *Manager) ExistsByUUID(uuid string) bool {
	m.Locker.RLock()
	defer m.Locker.RUnlock()
	if _, ok := m.Clients[uuid]; ok {
		return true
	}
	return false
}

// GetClientDomain 获取指定域名下的随机一个机器
func (m *Manager) GetClientByDomain(domain string) (*Client, bool) {
	m.Locker.RLock()
	defer m.Locker.RUnlock()

	if m.ExistsByDomain(domain) {
		return m.Domains[domain][rand.Intn(len(m.Domains[domain]))], true
	}

	return nil, false
}

// GetClientDomain 获取指定域名下所有机器
func (m *Manager) GetClientsByDomain(domain string) ([]*Client, bool) {
	m.Locker.RLock()
	defer m.Locker.RUnlock()

	if m.ExistsByDomain(domain) {
		return m.Domains[domain], true
	}

	return nil, false
}

// GetClientDomain 获取指定 uuid 的一个机器
func (m *Manager) GetClientByUUID(uuid string) (*Client, bool) {
	m.Locker.RLock()
	defer m.Locker.RUnlock()
	if v, ok := m.Clients[uuid]; ok {
		return v, true
	}
	return nil, false
}

// AddClient 添加一个机器
func (m *Manager) AddClient(conn *websocket.Conn, clientUUID string) *Client {
	m.Locker.Lock()
	defer m.Locker.Unlock()

	client := &Client{
		Conn: conn,
		UUID: clientUUID,
	}

	m.Clients[client.UUID] = client

	return client
}

// RemoveClient 移除一个机器
func (m *Manager) RemoveClient(uuid string) {
	client, exists := m.GetClientByUUID(uuid)

	if exists {
		delete(m.Clients, uuid)
		clients := m.Domains[client.Domain]
		for index, v := range clients {
			if v == client {
				m.Domains[client.Domain] = append(m.Domains[client.Domain][:index], m.Domains[client.Domain][index+1:]...)
			}
		}
	}
}

// ClientLength 获取机器的数量
func (m *Manager) ClientLength() int {
	return len(m.Clients)
}

// DomainLength 获取指定域名下的机器数量
func (m *Manager) DomainLength(domain string) int {
	if m.ExistsByDomain(domain) {
		return len(m.Domains[domain])
	}

	return 0
}

// CreateUUID 创建 UUID
func CreateUUID() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}
