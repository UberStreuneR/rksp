package server

import (
	"context"
	"fmt"
	"log"
	"practice04/services"
	"strconv"
	"strings"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"gorm.io/gorm"
)

type RsocketServer struct {
	m services.MessageService
	c services.ChannelService
	u services.UserService
}

func NewRsocketServer(db *gorm.DB) RsocketServer {
	m := services.CreateMessageService(db)
	c := services.CreateChannelService(db)
	u := services.CreateUserService(db)
	return RsocketServer{m, c, u}
}

// add channel
func (r RsocketServer) Log() func(payload.Payload) {
	return func(msg payload.Payload) {
		channelName := string(msg.Data())
		channel, err := r.c.AddOne(channelName)
		if err != nil {
			fmt.Println("Couldn't create channel: ", err)
			return
		}
		fmt.Println("Created channel: ", channel.Name)
	}
}

func (r RsocketServer) RequestResponse() func(payload.Payload) mono.Mono {
	return func(msg payload.Payload) mono.Mono {
		id, ok := msg.Metadata()
		if !ok {
			return mono.Just(payload.NewString("Couldn't parse id from msg metadata", ""))
		}
		strId, err := strconv.Atoi(string(id))
		if err != nil {
			return mono.Just(payload.NewString("Couldn't parse id from message strId", ""))
		}
		userID := uint(strId)
		data := strings.Split(string(msg.Data()), ":")
		channelName := data[0]
		message := data[1]
		res, err := r.m.AddOne(channelName, userID, message)
		if err != nil {
			log.Println(err)
			return mono.Just(payload.NewString("Couldn't create new msg", ""))
		}
		return mono.Just(payload.NewString(res.Content, ""))
	}
}

func (r RsocketServer) RequestStream() func(payload.Payload) flux.Flux {
	return func(msg payload.Payload) flux.Flux {
		channelName := string(msg.Data())
		msgs, err := r.m.GetAllFromChannel(channelName)
		return flux.Create(func(ctx context.Context, s flux.Sink) {
			if err != nil {
				return
			}
			for _, m := range msgs {
				newMsg := payload.NewString(fmt.Sprintf("%v: %v", m.User.ID, m.Content), "")
				s.Next(newMsg)
			}
			s.Complete()
		})
	}
}

func (r RsocketServer) RequestChannel() func(flux.Flux) flux.Flux {
	return func(requests flux.Flux) flux.Flux {
		requests = requests.Map(func(p payload.Payload) (payload.Payload, error) {
			return payload.NewString(reverseString(string(p.Data())), ""), nil
		})
		return requests
	}
}

func (r RsocketServer) Serve() {
	err := rsocket.Receive().
		Acceptor(func(ctx context.Context, setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			// Handle close.
			fmt.Println(string(setup.Data()))
			sendingSocket.OnClose(func(err error) {
				log.Println("sending socket is closed")
			})

			return rsocket.NewAbstractSocket(
				// rsocket.FireAndForget(func(msg payload.Payload) {
				// 	r.cur = string(msg.Data())
				// 	fmt.Println("Server -- FNF: set cur string to ", string(msg.Data()))
				// }),
				rsocket.FireAndForget(r.Log()),
				// rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
				// 	return mono.Just(payload.NewString(r.cur, ""))
				// }),
				rsocket.RequestResponse(r.RequestResponse()),
				// rsocket.RequestStream(func(msg payload.Payload) flux.Flux {
				// 	fmt.Println("Server - Request stream payload:", string(msg.Data()))
				// 	parts, err := strconv.Atoi(string(msg.Data()))
				// 	if err != nil {
				// 		panic("Couldn't parse int from RequestStream payload")
				// 	}
				// 	return flux.Create(func(ctx context.Context, s flux.Sink) {
				// 		for i := 0; i < parts; i++ {
				// 			ret := r.cur[i*len(r.cur)/parts : (i+1)*len(r.cur)/parts]
				// 			s.Next(payload.NewString(ret, fmt.Sprintf("This is response #%04d", i)))
				// 		}
				// 		s.Complete()
				// 	})
				// }),
				rsocket.RequestStream(r.RequestStream()),
				rsocket.RequestChannel(r.RequestChannel()),
				// rsocket.RequestChannel(func(requests flux.Flux) flux.Flux {
				// 	fmt.Println("Server - Request channel:")
				// 	requests = requests.Map(func(p payload.Payload) (payload.Payload, error) {
				// 		return payload.NewString(reverseString(string(p.Data())), ""), nil
				// 	})
				// 	return requests
				// }),
			), nil
		}).
		// Transport(rsocket.TCPServer().SetAddr(":7878").Build()).
		Transport(rsocket.WebsocketServer().SetAddr(":7878").Build()).
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
