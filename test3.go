package main

import (
	"context"
	"log"
	"fmt"
	"time"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"	
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
)

func rsocket_client2() {
	cli, err := rsocket.Connect().
	Resume(). // Enable RESUME.
	Lease().  // Enable LEASE.
	Fragment(4096).
	SetupPayload(payload.NewString("Hello", "World")).
	Acceptor(func(ctx context.Context, socket rsocket.RSocket) rsocket.RSocket {
		return rsocket.NewAbstractSocket(
			rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
				return mono.Just(payload.NewString("Pong", time.Now().String()))
			}),
		)
	}).
	Transport(rsocket.TCPClient().SetAddr("127.0.0.1:7878").Build()).
	Start(context.Background())
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = cli.Close()
	}()
	// Simple FireAndForget.
	cli.FireAndForget(payload.NewString("This is a FNF message.", ""))
	// Simple RequestResponse.
	cli.RequestResponse(payload.NewString("This is a RequestResponse message.", "")).
		DoOnSuccess(func(elem payload.Payload) error {
			log.Println("response:", elem)
			return nil
		}).
		Subscribe(context.Background())
	var s rx.Subscription
	// RequestStream with backpressure. (one by one)
	cli.RequestStream(payload.NewString("This is a RequestStream message.", "")).
		DoOnNext(func(elem payload.Payload) error {
			log.Println("next element in stream:", elem)
			s.Request(1)
			return nil
		}).
		Subscribe(context.Background(), rx.OnSubscribe(func(ctx context.Context, s rx.Subscription) {
			s.Request(1)
		}))
	// Simple RequestChannel.
	sendFlux := flux.Create(func(ctx context.Context, s flux.Sink) {
		for i := 0; i < 3; i++ {
			s.Next(payload.NewString(fmt.Sprintf("This is a RequestChannel message #%d.", i), ""))
		}
		s.Complete()
	})
	cli.RequestChannel(sendFlux).
		DoOnNext(func(elem payload.Payload) error {
			log.Println("next element in channel:", elem)
			return nil
		}).
		Subscribe(context.Background())

}

