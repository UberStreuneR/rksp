package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

func rsocket_client2() {
	// time.Sleep(time.Second)
	cli, err := rsocket.Connect().
		Fragment(4096).
		// SetupPayload(payload.NewString("Hello", "World")).
		Acceptor(func(ctx context.Context, socket rsocket.RSocket) rsocket.RSocket {
			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					return mono.Just(payload.NewString("Pong", time.Now().String()))
				}),
			)
		}).
		// Transport(rsocket.TCPClient().SetAddr("127.0.0.1:7878").Build()).
		Transport(rsocket.TCPClient().SetAddr("localhost:7878").Build()).
		Start(context.Background())
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = cli.Close()
	}()

	// Simple FireAndForget.
	cli.FireAndForget(payload.NewString("This will be a new current string for the server", "Metadata"))
	fmt.Println()

	// Simple RequestResponse.
	cli.RequestResponse(payload.NewString("This is a RequestResponse message.", "")).
		DoOnSuccess(func(elem payload.Payload) error {
			log.Println("Client -- RequestResponse: current string is", string(elem.Data()))
			return nil
		}).
		Subscribe(context.Background())
	// time.Sleep(time.Second)
	fmt.Println()

	var sub rx.Subscription
	var wg sync.WaitGroup
	wg.Add(1)
	// RequestStream with backpressure. (one by one)
	// split cur string into pieces of 5
	cli.RequestStream(payload.NewString("5", "")).
		DoOnNext(func(elem payload.Payload) error {
			log.Println("Client -- RequestStream: ", string(elem.Data()))
			sub.Request(1)
			return nil
		}).DoOnComplete(wg.Done).
		Subscribe(context.Background(), rx.OnSubscribe(func(ctx context.Context, s rx.Subscription) {
			sub = s
			s.Request(1)
		}))
	wg.Wait()
	fmt.Println()

	// // Simple RequestChannel.
	sendFlux := flux.Create(func(ctx context.Context, s flux.Sink) {
		str := "This is the string that we split and reverse"
		for i := 0; i < 3; i++ {
			cur := str[i*len(str)/3 : (i+1)*len(str)/3]
			s.Next(payload.NewString(cur, ""))
		}
		s.Complete()
	})
	wg.Add(1)
	cli.RequestChannel(sendFlux).
		DoOnNext(func(elem payload.Payload) error {
			log.Println("next element in channel:", string(elem.Data()))
			return nil
		}).DoOnComplete(wg.Done).
		Subscribe(context.Background())
	wg.Wait()
}
