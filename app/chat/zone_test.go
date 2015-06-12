package chat

import (
	"github.com/jpcummins/geochat/app/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ZoneTestSuite struct {
	suite.Suite
	cache *cache.MockCache
	world *World
	root  *Zone
}

func (suite *ZoneTestSuite) SetupTest() {
	suite.cache = &cache.MockCache{}
	suite.cache.On("SetZone", mock.Anything).Return(nil)

	world, err := newWorld(suite.cache, 2) // SetZone called for cache
	assert.NoError(suite.T(), err)
	suite.world = world
	suite.root = world.root
	suite.cache.AssertCalled(suite.T(), "SetZone", suite.world.root)
}

func (suite *ZoneTestSuite) TestNewZone() {
	assert.Equal(suite.T(), ":0z", suite.root.ID())
	assert.Equal(suite.T(), 90, suite.root.boundary.NorthEastLat)
	assert.Equal(suite.T(), 180, suite.root.boundary.NorthEastLong)
	assert.Equal(suite.T(), -90, suite.root.boundary.SouthWestLat)
	assert.Equal(suite.T(), -180, suite.root.boundary.SouthWestLong)
	assert.Equal(suite.T(), "", suite.root.geohash)
	assert.Equal(suite.T(), byte('0'), suite.root.from)
	assert.Equal(suite.T(), byte('z'), suite.root.to)
	assert.Nil(suite.T(), suite.root.parent)
	assert.Nil(suite.T(), suite.root.left)
	assert.Nil(suite.T(), suite.root.right)
	assert.Equal(suite.T(), 0, suite.root.Count())
	assert.True(suite.T(), suite.root.IsOpen())
}

func (suite *ZoneTestSuite) TestAddUser() {
	user := &mockUser{}
	user.On("ID").Return("user1")
	suite.root.AddUser(user)
	assert.Equal(suite.T(), 1, suite.root.Count())
	assert.True(suite.T(), suite.root.IsOpen())
	suite.cache.AssertNumberOfCalls(suite.T(), "SetZone", 2)
}

func (suite *ZoneTestSuite) TestMaxUsersClosesRoom() {
	user1 := &mockUser{}
	user1.On("ID").Return("user1")
	user2 := &mockUser{}
	user2.On("ID").Return("user2")
	suite.root.AddUser(user1) // SetZone called for cache
	suite.root.AddUser(user2) // SetZone called for cache
	assert.Equal(suite.T(), 2, suite.root.Count())
	assert.False(suite.T(), suite.root.IsOpen())
	suite.cache.AssertNumberOfCalls(suite.T(), "SetZone", 3)
}

func (suite *ZoneTestSuite) TestRemoveUser() {
	user1 := &mockUser{}
	user1.On("ID").Return("user1")
	assert.Equal(suite.T(), 0, suite.root.Count())
	suite.root.AddUser(user1)
	assert.Equal(suite.T(), 1, suite.root.Count())
	suite.root.RemoveUser("user1") // SetZone called
	assert.Equal(suite.T(), 0, suite.root.Count())
	suite.cache.AssertNumberOfCalls(suite.T(), "SetZone", 3)
}

func (suite *ZoneTestSuite) TestRemoveUserOpensRoom() {
	user1 := &mockUser{}
	user1.On("ID").Return("user1")
	user2 := &mockUser{}
	user2.On("ID").Return("user2")
	suite.root.AddUser(user1)
	suite.root.AddUser(user2)
	assert.False(suite.T(), suite.root.IsOpen())
	suite.root.RemoveUser("user1") // SetZone called
	assert.True(suite.T(), suite.root.IsOpen())
	assert.Equal(suite.T(), 1, suite.root.Count())
	suite.cache.AssertNumberOfCalls(suite.T(), "SetZone", 4)
}

func (suite *ZoneTestSuite) TestBroadcast() {
	event := &mockEvent{}

	user1 := &mockUser{}
	user1.On("ID").Return("user1")
	user1.On("Broadcast", event).Return(nil)

	user2 := &mockUser{}
	user2.On("ID").Return("user2")
	user2.On("Broadcast", event).Return(nil)

	suite.root.AddUser(user1)
	suite.root.AddUser(user2)
	suite.root.Broadcast(event)

	user1.AssertCalled(suite.T(), "Broadcast", event)
	user2.AssertCalled(suite.T(), "Broadcast", event)
}

func TestZoneTestSuit(t *testing.T) {
	suite.Run(t, new(ZoneTestSuite))
}
