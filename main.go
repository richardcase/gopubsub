package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	wg, ctx := errgroup.WithContext(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	publisher := &Publisher{}

	subscriber1 := NewSubscriber("subscriber1")
	publisher.Subscribe(subscriber1.Name, subscriber1.Messages)

	subscriber2 := NewSubscriber("subscriber2")
	publisher.Subscribe(subscriber2.Name, subscriber2.Messages)

	wg.Go(func() error {
		subscriber1.Run(ctx)
		fmt.Println("returning from subscriber1 go routine")
		return nil
	})
	wg.Go(func() error {
		subscriber2.Run(ctx)
		fmt.Println("returning from subscriber2 go routine")
		return nil
	})

	ticker := time.NewTicker(time.Second * 2)
	wg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("publisher go routine: shutdown received, shutdown....")
				return nil
			case t := <-ticker.C:
				message := fmt.Sprint("Message generated at ", t)
				publisher.Publish(message)
				fmt.Println("published message")
			}
		}
	})

	select {
	case <-signalChan:
		fmt.Println("Signal received: going to call cancel")
		cancel()
	case <-ctx.Done():
	}

	fmt.Println("going to call cancel")
	cancel()
	fmt.Println("going to wait for wait group")
	if err := wg.Wait(); err != nil {
		fmt.Printf("unhandled error, existing. %v", err)
	}

	fmt.Println("All done exiting")
}
