package services

import (
	"log"
	"practice04/initializers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserLogServiceTestSuite struct {
	suite.Suite
	u  UserService
	c  ChannelService
	ul UserLogService
}

func (suite *UserLogServiceTestSuite) TearDownSuite() {
	testDB, _ := suite.u.DB.DB()
	testDB.Close()
	initializers.DB.Exec("DROP DATABASE test")
}

func (suite *UserLogServiceTestSuite) TearDownTest() {
	// suite.u.DeleteAll()
	suite.ul.DeleteAll()
}

func (suite *UserLogServiceTestSuite) TestAddOne() {
	res, err := suite.ul.AddOne(1, "Test channel", "delete")
	suite.Equal(err, nil)
	suite.Equal(res.UserID, uint(1))
	suite.Equal(res.ChannelName, "Test channel")
	suite.Equal(res.Operation, "delete")
}

func (suite *UserLogServiceTestSuite) TestGetAll() {
	suite.ul.AddOne(1, "Test channel", "delete")
	suite.ul.AddOne(2, "Test channel", "delete")
	res, err := suite.ul.GetAll()
	suite.Equal(err, nil)
	suite.Equal(len(res), 2)
}

// func (suite *UserLogServiceTestSuite) TestAutoLog() {
// 	u, err := suite.u.AddOne(1)
// 	suite.Equal(err, nil)
// 	c1, _ := suite.c.AddOne("Test channel 1")
// 	c2, _ := suite.c.AddOne("Test channel 2")
// 	suite.c.AddChannelToUser(u.ID, []string{c1.Name, c2.Name})
// 	// logs, err := suite.sl.GetPeriod(u.ID, "2023-01", "2023-12")
// 	logs, err := suite.ul.GetAll()
// 	suite.Equal(err, nil)
// 	suite.Equal(len(logs), 2)
// 	suite.Equal(logs[0].UserID, u.ID)
// 	suite.Equal(logs[0].ChannelName, c1.Name)
// 	suite.Equal(logs[0].Operation, "added")
// 	suite.Equal(logs[1].UserID, u.ID)
// 	suite.Equal(logs[1].ChannelName, c2.Name)
// 	suite.Equal(logs[1].Operation, "added")

//		suite.c.RemoveUserFromChannel(u.ID, "Test channel 1")
//		suite.c.RemoveUserFromChannel(u.ID, "Test channel 2")
//		logs, err = suite.ul.GetAll()
//		suite.Equal(err, nil)
//		suite.Equal(len(logs), 4)
//		suite.Equal(logs[2].UserID, u.ID)
//		suite.Equal(logs[2].ChannelName, c1.Name)
//		suite.Equal(logs[2].Operation, "deleted")
//		suite.Equal(logs[3].UserID, u.ID)
//		suite.Equal(logs[3].ChannelName, c2.Name)
//		suite.Equal(logs[3].Operation, "deleted")
//	}
func TestSegmentLogService(t *testing.T) {
	config, err := initializers.LoadConfig("../")
	if err != nil {
		log.Fatalln("Failed to load environment variables\n", err.Error())
	}
	db := initializers.GetTestDB(&config)
	u := CreateUserService(db)
	c := CreateChannelService(db)
	ul := CreateUserLogService(db)
	suite.Run(t, &UserLogServiceTestSuite{u: u, c: c, ul: ul})
}
