package chat

import (
	"github.com/jpcummins/geochat/app/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
	cache *cache.MockUserCache
	world *World
	user  *User
}

func (suite *UserTestSuite) SetupTest() {
	suite.cache = &cache.MockUserCache{}
	suite.cache.On("Set", mock.Anything).Return(nil)
	suite.world = newWorld(suite.cache)

	u, err := suite.world.NewUser(47.6235616, -122.330341, "test", "testid")
	suite.user = u
	assert.NoError(suite.T(), err)
	suite.cache.AssertCalled(suite.T(), "Set", suite.user)
}

func (suite *UserTestSuite) TestNewUser() {
	assert.Equal(suite.T(), "test", suite.user.Name())
	assert.Equal(suite.T(), "testid", suite.user.ID())
	assert.Nil(suite.T(), suite.user.Zone())
}

func (suite *UserTestSuite) TestConnect() {
	connection, err := suite.user.NewConnection()
	assert.NotNil(suite.T(), connection.Events())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(suite.user.connections))
}

func (suite *UserTestSuite) TestDisconnect() {
	connection, _ := suite.user.NewConnection()
	suite.user.Disconnect(connection)
	assert.Equal(suite.T(), 0, len(suite.user.connections))
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
