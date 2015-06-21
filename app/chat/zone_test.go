package chat

import (
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ZoneTestSuite struct {
	suite.Suite
	world *mocks.World
}

func (suite *ZoneTestSuite) SetupTest() {
	suite.world = &mocks.World{}
	suite.world.On("ID").Return("worldid")
}

func (suite *ZoneTestSuite) TestNewZone() {
	zone, err := newZone(rootZoneID, suite.world, 2)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), rootZoneID, zone.ID())
	assert.Equal(suite.T(), float64(90), zone.NorthEast().Lat())
	assert.Equal(suite.T(), float64(180), zone.NorthEast().Lng())
	assert.Equal(suite.T(), float64(-90), zone.SouthWest().Lat())
	assert.Equal(suite.T(), float64(-180), zone.SouthWest().Lng())
	assert.Equal(suite.T(), "", zone.Geohash())
	assert.Equal(suite.T(), "0", zone.From())
	assert.Equal(suite.T(), "z", zone.To())
	// assert.Nil(suite.T(), zone.ParentZoneID())
	assert.Equal(suite.T(), ":0g", zone.LeftZoneID())
	assert.Equal(suite.T(), ":hz", zone.RightZoneID())
	assert.Equal(suite.T(), 0, zone.Count())
	assert.True(suite.T(), zone.IsOpen())
}

func (suite *ZoneTestSuite) TestAddUser() {
	user := &mocks.User{}
	user.On("ID").Return("user1")

	zone, err := newZone(rootZoneID, suite.world, 2)
	assert.NoError(suite.T(), err)

	zone.AddUser(user)
	assert.Equal(suite.T(), 1, zone.Count())
	assert.True(suite.T(), zone.IsOpen())
}

func (suite *ZoneTestSuite) TestSetIsOpen() {
	zone, err := newZone(rootZoneID, suite.world, 2)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), zone.IsOpen())
	zone.SetIsOpen(false)
	assert.False(suite.T(), zone.IsOpen())
	zone.SetIsOpen(true)
	assert.True(suite.T(), zone.IsOpen())
}

func (suite *ZoneTestSuite) TestRemoveUser() {
	user1 := &mocks.User{}
	user1.On("ID").Return("user1")

	zone, err := newZone(rootZoneID, suite.world, 2)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 0, zone.Count())
	zone.AddUser(user1)
	assert.Equal(suite.T(), 1, zone.Count())
	zone.RemoveUser("user1")
	assert.Equal(suite.T(), 0, zone.Count())
}

// func (suite *ZoneTestSuite) TestBroadcast() {
// 	event := &mocks.Event{}
//
// 	user1 := &mocks.User{}
// 	user1.On("ID").Return("user1")
// 	user1.On("Broadcast", event).Return(nil)
//
// 	user2 := &mocks.User{}
// 	user2.On("ID").Return("user2")
// 	user2.On("Broadcast", event).Return(nil)
//
// 	zone, err := newZone(rootZoneID, suite.world, 2)
// 	assert.NoError(suite.T(), err)
//
// 	zone.AddUser(user1)
// 	zone.AddUser(user2)
// 	zone.Broadcast(event)
//
// 	user1.AssertCalled(suite.T(), "Broadcast", event)
// 	user2.AssertCalled(suite.T(), "Broadcast", event)
// }

func (suite *ZoneTestSuite) TestMarshalJSON() {
	zone, err := newZone(rootZoneID, suite.world, 2)
	assert.NoError(suite.T(), err)

	user1 := &mocks.User{}
	user1.On("ID").Return("user1")
	zone.AddUser(user1)

	user2 := &mocks.User{}
	user2.On("ID").Return("user2")
	zone.AddUser(user2)

	b, err := zone.MarshalJSON()
	assert.Equal(suite.T(), "{\"id\":\":0z\",\"user_ids\":[\"user1\",\"user2\"],\"is_open\":true,\"max_users\":2}", string(b))
	assert.NoError(suite.T(), err)
}

func TestZoneTestSuit(t *testing.T) {
	suite.Run(t, new(ZoneTestSuite))
}
