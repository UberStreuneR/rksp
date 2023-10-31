package services

import (
	"log"
	"practice04/initializers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MessageServiceTestSuite struct {
	suite.Suite
	u  UserService
	c  ChannelService
	ul UserLogService
	m  MessageService
}

func (suite *MessageServiceTestSuite) TearDownSuite() {
	testDB, _ := suite.u.DB.DB()
	testDB.Close()
	initializers.DB.Exec("DROP DATABASE test")
}

func (suite *MessageServiceTestSuite) TearDownTest() {
	suite.c.DeleteAll()
	suite.u.DeleteAll()
	suite.ul.DeleteAll()
	suite.m.DeleteAll()
}

func (suite *MessageServiceTestSuite) TestAddOne() {
	channel, _ := suite.c.AddOne("Test channel")
	user, _ := suite.u.AddOne(1)
	suite.c.AddChannelToUser(user.ID, []string{channel.Name})
	res, err := suite.m.AddOne(channel.Name, user.ID, "rofl")
	suite.Equal(err, nil)
	suite.Equal(res.ID, uint(1))
	suite.Equal(res.UserID, uint(1))
	suite.Equal(res.Channel.Name, "Test channel")

}

func (suite *MessageServiceTestSuite) TestGetAllFromChannel() {
	channel, _ := suite.c.AddOne("Test channel")
	user, _ := suite.u.AddOne(1)
	user2, _ := suite.u.AddOne(2)
	suite.c.AddChannelToUser(user.ID, []string{channel.Name})
	suite.c.AddChannelToUser(user2.ID, []string{channel.Name})
	suite.m.AddOne(channel.Name, user.ID, "rofl")
	suite.m.AddOne(channel.Name, user.ID, "kekw")
	suite.m.AddOne(channel.Name, user.ID, "cringe")
	suite.m.AddOne(channel.Name, user2.ID, "hello")
	suite.m.AddOne(channel.Name, user2.ID, "why you spamming bro")
	messages, err := suite.m.GetAllFromChannel(channel.Name)
	suite.Equal(err, nil)
	suite.Equal(len(messages), 5)
	suite.Equal(messages[4].Content, "why you spamming bro")
}

func TestMessageService(t *testing.T) {
	config, err := initializers.LoadConfig("../")
	if err != nil {
		log.Fatalln("Failed to load environment variables\n", err.Error())
	}
	db := initializers.GetTestDB(&config)
	u := CreateUserService(db)
	c := CreateChannelService(db)
	ul := CreateUserLogService(db)
	m := CreateMessageService(db)
	suite.Run(t, &MessageServiceTestSuite{u: u, c: c, ul: ul, m: m})
}
