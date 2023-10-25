package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

type RsocketServer struct {
	cur string
}

func (r RsocketServer) serve() {
	err := rsocket.Receive().
		Acceptor(func(ctx context.Context, setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			// Handle close.
			fmt.Println(string(setup.Data()))
			sendingSocket.OnClose(func(err error) {
				log.Println("sending socket is closed")
			})

			return rsocket.NewAbstractSocket(
				rsocket.FireAndForget(func(msg payload.Payload) {
					r.cur = string(msg.Data())
					fmt.Println("Server -- FNF: set cur string to ", string(msg.Data()))
				}),
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					return mono.Just(payload.NewString(r.cur, ""))
				}),
				rsocket.RequestStream(func(msg payload.Payload) flux.Flux {
					fmt.Println("Server - Request stream payload:", string(msg.Data()))
					parts, err := strconv.Atoi(string(msg.Data()))
					if err != nil {
						panic("Couldn't parse int from RequestStream payload")
					}
					return flux.Create(func(ctx context.Context, s flux.Sink) {
						for i := 0; i < parts; i++ {
							ret := r.cur[i*len(r.cur)/parts : (i+1)*len(r.cur)/parts]
							s.Next(payload.NewString(ret, fmt.Sprintf("This is response #%04d", i)))
						}
						s.Complete()
					})
				}),
				rsocket.RequestChannel(func(requests flux.Flux) flux.Flux {
					fmt.Println("Server - Request channel:")
					requests = requests.Map(func(p payload.Payload) (payload.Payload, error) {
						return payload.NewString(reverseString(string(p.Data())), ""), nil
					})
					return requests
				}),
			), nil
		}).
		Transport(rsocket.TCPServer().SetAddr(":7878").Build()).
		Serve(context.Background())
	log.Fatalln("err:", err)
}

func reverseString(s string) string {
	r := []rune(s)
	for i := 0; i < len(s)/2; i++ {
		r[i], r[len(s)-i-1] = r[len(s)-i-1], r[i]
	}
	return string(r)
}
