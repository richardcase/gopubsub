package main

import (
	"context"
	"fmt"
)

type Subscriber struct {
	Name     string
	Messages chan interface{}
}

func NewSubscriber(name string) *Subscriber {
	return &Subscriber{
		Name:     name,
		Messages: make(chan interface{}),
	}
}

func (s *Subscriber) Run(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: shutdown received, shutdown....\n", s.Name)
			return
		case message := <-s.Messages:
			fmt.Printf("%s: %s\n", s.Name, message)
		}
	}
}
