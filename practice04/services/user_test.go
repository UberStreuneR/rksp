package services

import (
	"log"
	"practice04/initializers"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	u UserService
}

func (suite *UserServiceTestSuite) TestAddUser() {
	user, err := suite.u.AddOne(1000)
	suite.Equal(err, nil)
	userInDB, err := suite.u.GetOne(1000)
	suite.Equal(err, nil)
	suite.Equal(user.ID, userInDB.ID)
	suite.Equal(len(user.Channels), len(user.Channels))
}

func (suite *UserServiceTestSuite) TestGetAllUsers() {
	suite.u.AddOne(1001)
	suite.u.AddOne(1002)
	users, _ := suite.u.GetAll()
	suite.Equal(len(users), 3)
	suite.Equal(users[0].ID, uint(1000))
	suite.Equal(users[1].ID, uint(1001))
	suite.Equal(users[2].ID, uint(1002))
}

func (suite *UserServiceTestSuite) TearDownSuite() {
	testDB, _ := suite.u.DB.DB()
	testDB.Close()
	initializers.DB.Exec("DROP DATABASE test")
}

func TestUserService(t *testing.T) {
	config, err := initializers.LoadConfig("../")
	if err != nil {
		log.Fatalln("Failed to load environment variables\n", err.Error())
	}
	db := initializers.GetTestDB(&config)
	suite.Run(t, &UserServiceTestSuite{u: CreateUserService(db)})
}
