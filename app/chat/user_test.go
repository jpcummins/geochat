package chat

import (
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
	user *User
}

func (suite *UserTestSuite) SetupTest() {
	suite.user = NewUser("testid", "test", seattle)
}

func (suite *UserTestSuite) TestNewUser() {
	assert.Equal(suite.T(), "test", suite.user.Name())
	assert.Equal(suite.T(), "testid", suite.user.ID())
}

func (suite *UserTestSuite) TestAddConnection() {
	connection := &mocks.Connection{}
	suite.user.AddConnection(connection)
	assert.Equal(suite.T(), 1, len(suite.user.connections))
}

func (suite *UserTestSuite) TestMultipleConnections() {
	suite.user.AddConnection(&mocks.Connection{})
	suite.user.AddConnection(&mocks.Connection{})
	assert.Equal(suite.T(), 2, len(suite.user.connections))
}

func (suite *UserTestSuite) TestDisconnect() {
	c1 := &mocks.Connection{}
	suite.user.AddConnection(c1)
	assert.Equal(suite.T(), 1, len(suite.user.connections))
	suite.user.RemoveConnection(c1)
	assert.Equal(suite.T(), 0, len(suite.user.connections))
}

func (suite *UserTestSuite) TestDisconnectWithMultipleConnections() {
	connection1 := &mocks.Connection{}
	connection2 := &mocks.Connection{}
	connection3 := &mocks.Connection{}

	suite.user.AddConnection(connection1)
	suite.user.AddConnection(connection2)
	suite.user.AddConnection(connection3)

	assert.Equal(suite.T(), 3, len(suite.user.connections))
	suite.user.RemoveConnection(connection2)
	assert.Equal(suite.T(), 2, len(suite.user.connections))
	assert.Equal(suite.T(), connection1, suite.user.connections[0])
	assert.Equal(suite.T(), connection3, suite.user.connections[1])
}

func (suite *UserTestSuite) TestBroadcast() {
	connection1 := &mocks.Connection{}
	ch1 := make(chan types.Event, 1)
	connection1.On("Events").Return(ch1)

	connection2 := &mocks.Connection{}
	ch2 := make(chan types.Event, 1)
	connection2.On("Events").Return(ch2)

	connection3 := &mocks.Connection{}
	ch3 := make(chan types.Event, 1)
	connection3.On("Events").Return(ch3)

	suite.user.AddConnection(connection1)
	suite.user.AddConnection(connection2)
	suite.user.AddConnection(connection3)

	event := &mocks.Event{}
	suite.user.Broadcast(event)

	connection1.AssertCalled(suite.T(), "Events")
	connection2.AssertCalled(suite.T(), "Events")
	connection3.AssertCalled(suite.T(), "Events")

	assert.Equal(suite.T(), event, <-ch1)
	assert.Equal(suite.T(), event, <-ch2)
	assert.Equal(suite.T(), event, <-ch3)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
