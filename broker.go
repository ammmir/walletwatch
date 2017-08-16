package walletwatch

import (
	"log"
	"sync/atomic"
	"time"
)

func NewBroker() *Broker {
	return &Broker{
		IncomingTx:     make(chan Transaction),
		NewClients:     make(chan BrokerClientConfig),
		ClosingClients: make(chan BrokerClientConfig),
		clients:        make(map[uint64]BrokerClientConfig),
	}
}

type Broker struct {
	counter        uint64
	IncomingTx     chan Transaction
	NewClients     chan BrokerClientConfig
	ClosingClients chan BrokerClientConfig
	clients        map[uint64]BrokerClientConfig
}

func (t *Broker) Start() {
	for {
		select {
		case tx := <-t.IncomingTx:
			// deliver transactions to interested clients
			for _, cfg := range t.clients {
				var ok bool

				for addr, _ := range tx.Outputs() {
					if cfg.Addresses[addr] {
						ok = true
						break
					}
				}

				if ok {
					select {
					case cfg.Ch <- tx:
					case <-time.After(1 * time.Second):
						// TODO: slow client, drop it
						log.Println("slow client")
					}
				}
			}
		case cfg := <-t.NewClients:
			t.clients[cfg.Id] = cfg
		case ch := <-t.ClosingClients:
			delete(t.clients, ch.Id)
		}
	}
}

func (t *Broker) Id() uint64 {
	atomic.AddUint64(&t.counter, 1) // FIXME: lame, but lock-free and easy
	i := atomic.LoadUint64(&t.counter)
	return i
}

type BrokerClientConfig struct {
	Id        uint64
	Ch        chan Transaction
	Addresses map[string]bool
}
