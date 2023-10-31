package server

import (
	"context"
	"fmt"
	"log"
	"practice04/initializers"
	"testing"
	"time"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/mono"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	server RsocketServer
	client rsocket.Client
}

func (suite *ServerTestSuite) TearDownSuite() {
	testDB, _ := suite.server.u.DB.DB()
	testDB.Close()
	initializers.DB.Exec("DROP DATABASE test")
}

func (suite *ServerTestSuite) TearDownTest() {
	suite.server.c.DeleteAll()
	suite.server.u.DeleteAll()
	suite.server.m.DeleteAll()
}

func (suite *ServerTestSuite) AddChannel() {
	suite.NotNil(suite.client)
	suite.client.FireAndForget(payload.NewString("Test channel one", ""))
	suite.client.FireAndForget(payload.NewString("Test channel two", ""))
	time.Sleep(2)
	channels, err := suite.server.c.GetAll()
	suite.Equal(err, nil)
	suite.Equal(len(channels), 2)
	suite.Equal(channels[0].Name, "Test channel one")
	suite.Equal(channels[1].Name, "Test channel two")
}

func (suite *ServerTestSuite) AddMessage() {
	suite.server.u.AddOne(uint(1))
	suite.client.RequestResponse(payload.NewString("Test channel one:Message from user 1", "1")).
		DoOnSuccess(func(elem payload.Payload) error {
			suite.Equal(string(elem.Data()), "Message from user 1")
			messages, err := suite.server.m.GetAllFromChannel("Test channel one")
			suite.Equal(err, nil)
			suite.Equal(len(messages), 1)
			suite.Equal(messages[0].Content, "Message from user 2")
			suite.Equal(messages[0].UserID, uint(1))
			return nil
		}).
		Subscribe(context.Background())
}

func (suite *ServerTestSuite) RequestMessagesStream() {
	suite.server.u.AddOne(uint(1))
	suite.server.c.AddOne("Test channel 1")
	messages := []string{"Hello", "World", "It's me"}
	for _, m := range messages {
		suite.server.m.AddOne("Test channel 1", uint(1), m)
	}
	var sub rx.Subscription
	count := 0
	suite.client.RequestStream(payload.NewString("Test channel 1", "")).
		DoOnNext(func(elem payload.Payload) error {
			suite.Equal(string(elem.Data()), messages[count])
			sub.Request(1)
			count++
			return nil
		}).
		Subscribe(context.Background(), rx.OnSubscribe(func(ctx context.Context, s rx.Subscription) {
			sub = s
		}))
}

func Test_Rsocket_Server(t *testing.T) {
	config, err := initializers.LoadConfig("../")
	if err != nil {
		log.Fatalln("Failed to load environment variables\n", err.Error())
	}
	db := initializers.GetTestDB(&config)
	server := NewRsocketServer(db)
	go server.Serve()
	client, err := rsocket.Connect().
		Fragment(4096).
		// SetupPayload(payload.NewString("Hello", "World")).
		Acceptor(func(ctx context.Context, socket rsocket.RSocket) rsocket.RSocket {
			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					return mono.Just(payload.NewString("Pong", time.Now().String()))
				}),
			)
		}).
		// Transport(rsocket.TCPClient().SetAddr("localhost:7878").Build()).
		Transport(rsocket.WebsocketClient().SetURL("ws://localhost:7878").Build()).
		Start(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = client.Close()
	}()
	fmt.Println("Client: ", client)
	suite.Run(t, &ServerTestSuite{server: NewRsocketServer(db), client: client})
}
