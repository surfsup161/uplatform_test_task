package memcached

import (
	"log"
	"net"
	"sync"
)

type Memcached struct {
	mu        sync.Mutex
	conns     []net.Conn
	connsChan chan net.Conn
}

func New() *Memcached {
	return &Memcached{}
}

func (m *Memcached) OpenConnections(host string, count int) error {
	if count == 0 {
		count = 1
	}

	m.connsChan = make(chan net.Conn, count)
	for i := 0; i < count; i++ {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			return err
		}

		m.mu.Lock()
		m.conns = append(m.conns, conn)
		m.mu.Unlock()

		m.connsChan <- conn
	}

	return nil
}

func (m *Memcached) GetConn() net.Conn {
	return <-m.connsChan
}

func (m *Memcached) CloseConnections() {
	m.mu.Lock()
	for i := range m.conns {
		err := m.conns[i].Close()
		if err != nil {
			log.Println(err)
		}
	}
	m.mu.Unlock()
}
