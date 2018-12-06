package main

import (
	"sync"
	"time"
)

type Publisher struct {
	interval    time.Duration
	subLock     sync.Mutex
	subscribers map[string]chan<- interface{}
}

func NewPublisher(interval time.Duration) *Publisher {
	return &Publisher{
		interval: interval,
	}
}

func (p *Publisher) Publish(message interface{}) error {
	for _, subChan := range p.subscribers {
		subChan <- message
	}

	return nil
}

func (p *Publisher) Subscribe(name string, recvChan chan<- interface{}) {
	p.subLock.Lock()
	defer p.subLock.Unlock()

	if p.subscribers == nil {
		p.subscribers = make(map[string]chan<- interface{})
	}
	p.subscribers[name] = recvChan
}

func (p *Publisher) Unsubscribe(name string) {
	p.subLock.Lock()
	defer p.subLock.Unlock()

	delete(p.subscribers, name)

}
