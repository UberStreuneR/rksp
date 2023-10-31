package services

import (
	"log"
	"practice04/initializers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ChannelServiceTestSuite struct {
	suite.Suite
	u  UserService
	c  ChannelService
	ul UserLogService
}

func (suite *ChannelServiceTestSuite) TearDownSuite() {
	testDB, _ := suite.u.DB.DB()
	testDB.Close()
	initializers.DB.Exec("DROP DATABASE test")
}

func (suite *ChannelServiceTestSuite) TearDownTest() {
	suite.c.DeleteAll()
	suite.u.DeleteAll()
	suite.ul.DeleteAll()
}

func (suite *ChannelServiceTestSuite) TestAddOne() {
	res, err := suite.c.AddOne("Test channel")
	suite.Equal(err, nil)
	suite.Equal(res.ID, uint(1))
	suite.Equal(res.Name, "Test channel")
	suite.Equal(len(res.Users), 0)
}

func (suite *ChannelServiceTestSuite) TestGetAll() {
	suite.c.AddOne("Test channel 1")
	suite.c.AddOne("Test channel 2")
	res, err := suite.c.GetAll()
	suite.Equal(err, nil)
	suite.Equal(len(res), 2)
}

func (suite *ChannelServiceTestSuite) TestAddUserToChannel() {
	channel, _ := suite.c.AddOne("Test channel 1")
	user, _ := suite.u.AddOne(2)
	err := suite.c.AddChannelToUser(user.ID, []string{channel.Name})
	suite.Equal(err, nil)

	user, err = suite.u.GetOne(2)
	suite.Equal(user.ID, uint(2))
	suite.Equal(err, nil)
	suite.Equal(len(user.Channels), 1)
	channel, _ = suite.c.GetOne("Test channel 1")
	suite.Equal(len(channel.Users), 1)

	err = suite.c.RemoveUserFromChannel(user.ID, "Test channel 1")
	suite.Equal(err, nil)

	user, _ = suite.u.GetOne(1)
	suite.Equal(len(user.Channels), 0)
	channel, _ = suite.c.GetOne("Test channel 1")
	suite.Equal(len(channel.Users), 0)
}

func TestChannelService(t *testing.T) {
	config, err := initializers.LoadConfig("../")
	if err != nil {
		log.Fatalln("Failed to load environment variables\n", err.Error())
	}
	db := initializers.GetTestDB(&config)
	u := CreateUserService(db)
	c := CreateChannelService(db)
	ul := CreateUserLogService(db)
	suite.Run(t, &ChannelServiceTestSuite{u: u, c: c, ul: ul})
}
