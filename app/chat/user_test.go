package chat

import (
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
	world *mocks.World
	user  *User
}

func (suite *UserTestSuite) SetupTest() {
	suite.world = &mocks.World{}
	suite.world.On("ID").Return("worldid")
	suite.user = newUser("testid", suite.world)
}

func (suite *UserTestSuite) TestNewUser() {
	assert.Equal(suite.T(), "testid", suite.user.ID())
}

func (suite *UserTestSuite) TestAddConnection() {
	suite.user.Connect()
	assert.Equal(suite.T(), 1, len(suite.user.connections))
}

func (suite *UserTestSuite) TestMultipleConnections() {
	suite.user.Connect()
	suite.user.Connect()
	assert.Equal(suite.T(), 2, len(suite.user.connections))
}

func (suite *UserTestSuite) TestDisconnect() {
	c1 := suite.user.Connect()
	assert.Equal(suite.T(), 1, len(suite.user.connections))
	suite.user.Disconnect(c1)
	assert.Equal(suite.T(), 0, len(suite.user.connections))
}

func (suite *UserTestSuite) TestDisconnectWithMultipleConnections() {
	connection1 := suite.user.Connect()
	connection2 := suite.user.Connect()
	connection3 := suite.user.Connect()
	assert.Equal(suite.T(), 3, len(suite.user.connections))

	suite.user.Disconnect(connection2)
	assert.Equal(suite.T(), 2, len(suite.user.connections))
	assert.Equal(suite.T(), connection1, suite.user.connections[0])
	assert.Equal(suite.T(), connection3, suite.user.connections[1])
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
