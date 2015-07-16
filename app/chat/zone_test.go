package chat

import (
	"encoding/json"
	"errors"
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ZoneTestSuite struct {
	suite.Suite
	world  *mocks.World
	logger *mocks.Logger
	zones  *mocks.Zones
	users  *mocks.Users
	db     *mocks.DB
}

func (suite *ZoneTestSuite) SetupTest() {
	suite.world = &mocks.World{}
	suite.logger = &mocks.Logger{}
	suite.zones = &mocks.Zones{}
	suite.users = &mocks.Users{}
	suite.db = &mocks.DB{}

	suite.world.On("ID").Return("worldid")
	suite.world.On("Zones").Return(suite.zones)
	suite.world.On("Users").Return(suite.users)
	suite.world.On("DB").Return(suite.db)
	suite.logger.On("New", mock.Anything).Return(suite.logger)
	suite.zones.On("UpdateCache", mock.Anything).Return(nil)
}

func (suite *ZoneTestSuite) TestNewZone() {
	zone, err := newZone(rootZoneID, suite.world, suite.logger)
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

func (suite *ZoneTestSuite) TestNewZone_Error() {
	_, err := newZone(invalidZoneID, suite.world, suite.logger)
	assert.Error(suite.T(), err)
}

func (suite *ZoneTestSuite) TestAddUser() {
	user := &mocks.User{}
	user.On("ID").Return("user1")

	zone, err := newZone(rootZoneID, suite.world, suite.logger)
	assert.NoError(suite.T(), err)

	zone.AddUser(user)
	assert.Equal(suite.T(), 1, zone.Count())
	assert.True(suite.T(), zone.IsOpen())
	assert.Equal(suite.T(), "user1", zone.UserIDs()[0])
}

func (suite *ZoneTestSuite) TestAddUser_Exists() {
	user := &mocks.User{}
	user.On("ID").Return("user1")

	zone, err := newZone(rootZoneID, suite.world, suite.logger)
	assert.NoError(suite.T(), err)

	zone.AddUser(user)
	assert.Equal(suite.T(), 1, zone.Count())
	zone.AddUser(user)
	assert.Equal(suite.T(), 1, zone.Count())
}

func (suite *ZoneTestSuite) TestSetIsOpen() {
	zone, err := newZone(rootZoneID, suite.world, suite.logger)
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

	zone, err := newZone(rootZoneID, suite.world, suite.logger)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 0, zone.Count())
	zone.AddUser(user1)
	assert.Equal(suite.T(), 1, zone.Count())
	zone.RemoveUser("user1")
	assert.Equal(suite.T(), 0, zone.Count())
}

func (suite *ZoneTestSuite) TestBroadcast() {
	event := &mocks.BroadcastEventData{}

	user1 := &mocks.User{}
	user1.On("ID").Return("user1")
	user1.On("Broadcast", event).Return(nil)

	user2 := &mocks.User{}
	user2.On("ID").Return("user2")
	user2.On("Broadcast", event).Return(nil)

	suite.users.On("User", "user1").Return(user1, nil)
	suite.users.On("User", "user2").Return(user2, nil)

	zone, err := newZone(rootZoneID, suite.world, suite.logger)
	assert.NoError(suite.T(), err)

	zone.AddUser(user1)
	zone.AddUser(user2)
	zone.Broadcast(event)

	user1.AssertCalled(suite.T(), "Broadcast", event)
	user2.AssertCalled(suite.T(), "Broadcast", event)
}

func (suite *ZoneTestSuite) TestBroadcastJSON() {
	zone, err := newZone(rootZoneID, suite.world, suite.logger)
	assert.NoError(suite.T(), err)

	user1 := &mocks.User{}
	user1.On("ID").Return("user1")
	user1.On("BroadcastJSON").Return(&types.UserBroadcastJSON{"user1", "User1", nil})
	zone.AddUser(user1)

	user2 := &mocks.User{}
	user2.On("ID").Return("user2")
	user2.On("BroadcastJSON").Return(&types.UserBroadcastJSON{"user2", "User2", nil})
	zone.AddUser(user2)

	suite.users.On("User", "user1").Return(user1, nil)
	suite.users.On("User", "user2").Return(user2, nil)

	b, err := json.Marshal(zone.BroadcastJSON())
	assert.Equal(suite.T(), "{\"id\":\":0z\",\"users\":{\"user1\":{\"id\":\"user1\",\"name\":\"User1\",\"location\":null},\"user2\":{\"id\":\"user2\",\"name\":\"User2\",\"location\":null}},\"sw\":{\"lat\":-90,\"lng\":-180},\"ne\":{\"lat\":90,\"lng\":180}}", string(b))
	assert.NoError(suite.T(), err)
}

func (suite *ZoneTestSuite) TestMarshalJSON() {
	zone, err := newZone(rootZoneID, suite.world, suite.logger)
	assert.NoError(suite.T(), err)

	user1 := &mocks.User{}
	user1.On("ID").Return("user1")
	zone.AddUser(user1)

	user2 := &mocks.User{}
	user2.On("ID").Return("user2")
	zone.AddUser(user2)

	b, err := json.Marshal(zone.PubSubJSON())
	assert.Equal(suite.T(), "{\"id\":\":0z\",\"user_ids\":[\"user1\",\"user2\"],\"is_open\":true,\"last_split_time\":0,\"last_merge_time\":0}", string(b))
	assert.NoError(suite.T(), err)
}

func (suite *ZoneTestSuite) TestUpdate() {
	zone := &Zone{}
	zone.world = suite.world
	suite.world.On("Zones").Return(suite.zones)
	suite.zones.On("UpdateCache", mock.Anything).Return(nil)

	update := &types.ZonePubSubJSON{
		ID:      ":0g",
		UserIDs: []string{"1", "2", "3"},
		IsOpen:  false,
	}

	zone.Update(update)

	assert.Equal(suite.T(), ":0g", zone.ID())
	assert.Equal(suite.T(), 3, zone.Count())
	assert.Equal(suite.T(), []string{"1", "2", "3"}, zone.UserIDs())
	assert.Equal(suite.T(), false, zone.IsOpen())
	suite.zones.AssertCalled(suite.T(), "UpdateCache", zone)
}

func (suite *ZoneTestSuite) TestUpdate_Error() {
	zone := &Zone{}
	err := zone.Update(&mocks.PubSubJSON{})
	assert.Equal(suite.T(), "Unable to serialize to ZonePubSubJSON.", err.Error())
}

/*
level 1 - 0123456789bcdefghjkmnpqrstuvwxyz
level 2 - 0123456789bcdefg hjkmnpqrstuvwxyz
level 3 - 01234567 89bcdefg hjkmnpqr stuvwxyz
level 4 - 0123 4567 89bc defg hjkm npqr stuv wxyz
level 5 - 01 23 45 67 89 bc de fg hj km np qr st uv wx yz

          1         2         3
01234567890123456789012345678901234567890
0123456789bcdefghjkmnpqrstuvwxyz
*/

func (suite *ZoneTestSuite) TestParentID_Level1() {
	zone, _ := newZone(":0z", suite.world, suite.logger)
	assert.Equal(suite.T(), "", zone.ParentZoneID())

	zone, _ = newZone("0:0z", suite.world, suite.logger)
	assert.Equal(suite.T(), ":01", zone.ParentZoneID())

	zone, _ = newZone("1:0z", suite.world, suite.logger)
	assert.Equal(suite.T(), ":01", zone.ParentZoneID())
}

func (suite *ZoneTestSuite) TestParentID_Level2() {
	zone, _ := newZone(":0g", suite.world, suite.logger)
	assert.Equal(suite.T(), ":0z", zone.ParentZoneID())

	zone, _ = newZone(":hz", suite.world, suite.logger)
	assert.Equal(suite.T(), ":0z", zone.ParentZoneID())
}

func (suite *ZoneTestSuite) TestParentID_Level3() {
	zone, _ := newZone(":07", suite.world, suite.logger)
	assert.Equal(suite.T(), ":0g", zone.ParentZoneID())

	zone, _ = newZone(":8g", suite.world, suite.logger)
	assert.Equal(suite.T(), ":0g", zone.ParentZoneID())

	zone, _ = newZone(":hr", suite.world, suite.logger)
	assert.Equal(suite.T(), ":hz", zone.ParentZoneID())

	zone, _ = newZone(":sz", suite.world, suite.logger)
	assert.Equal(suite.T(), ":hz", zone.ParentZoneID())
}

func (suite *ZoneTestSuite) TestParentID_Level4() {
	zone, _ := newZone(":03", suite.world, suite.logger)
	assert.Equal(suite.T(), ":07", zone.ParentZoneID())

	zone, _ = newZone(":47", suite.world, suite.logger)
	assert.Equal(suite.T(), ":07", zone.ParentZoneID())

	zone, _ = newZone(":8c", suite.world, suite.logger)
	assert.Equal(suite.T(), ":8g", zone.ParentZoneID())

	zone, _ = newZone(":dg", suite.world, suite.logger)
	assert.Equal(suite.T(), ":8g", zone.ParentZoneID())

	zone, _ = newZone(":hm", suite.world, suite.logger)
	assert.Equal(suite.T(), ":hr", zone.ParentZoneID())

	zone, _ = newZone(":nr", suite.world, suite.logger)
	assert.Equal(suite.T(), ":hr", zone.ParentZoneID())

	zone, _ = newZone(":sv", suite.world, suite.logger)
	assert.Equal(suite.T(), ":sz", zone.ParentZoneID())

	zone, _ = newZone(":wz", suite.world, suite.logger)
	assert.Equal(suite.T(), ":sz", zone.ParentZoneID())
}

func (suite *ZoneTestSuite) TestParentID_Level5() {
	zone, _ := newZone(":01", suite.world, suite.logger)
	assert.Equal(suite.T(), ":03", zone.ParentZoneID())

	zone, _ = newZone(":23", suite.world, suite.logger)
	assert.Equal(suite.T(), ":03", zone.ParentZoneID())

	zone, _ = newZone(":45", suite.world, suite.logger)
	assert.Equal(suite.T(), ":47", zone.ParentZoneID())

	zone, _ = newZone(":67", suite.world, suite.logger)
	assert.Equal(suite.T(), ":47", zone.ParentZoneID())

	zone, _ = newZone(":89", suite.world, suite.logger)
	assert.Equal(suite.T(), ":8c", zone.ParentZoneID())

	zone, _ = newZone(":bc", suite.world, suite.logger)
	assert.Equal(suite.T(), ":8c", zone.ParentZoneID())

	zone, _ = newZone(":de", suite.world, suite.logger)
	assert.Equal(suite.T(), ":dg", zone.ParentZoneID())

	zone, _ = newZone(":fg", suite.world, suite.logger)
	assert.Equal(suite.T(), ":dg", zone.ParentZoneID())

	zone, _ = newZone(":hj", suite.world, suite.logger)
	assert.Equal(suite.T(), ":hm", zone.ParentZoneID())

	zone, _ = newZone(":km", suite.world, suite.logger)
	assert.Equal(suite.T(), ":hm", zone.ParentZoneID())

	zone, _ = newZone(":np", suite.world, suite.logger)
	assert.Equal(suite.T(), ":nr", zone.ParentZoneID())

	zone, _ = newZone(":qr", suite.world, suite.logger)
	assert.Equal(suite.T(), ":nr", zone.ParentZoneID())

	zone, _ = newZone(":st", suite.world, suite.logger)
	assert.Equal(suite.T(), ":sv", zone.ParentZoneID())

	zone, _ = newZone(":uv", suite.world, suite.logger)
	assert.Equal(suite.T(), ":sv", zone.ParentZoneID())

	zone, _ = newZone(":wx", suite.world, suite.logger)
	assert.Equal(suite.T(), ":wz", zone.ParentZoneID())

	zone, _ = newZone(":yz", suite.world, suite.logger)
	assert.Equal(suite.T(), ":wz", zone.ParentZoneID())
}

func (suite *ZoneTestSuite) TestJoin_LeaveIsCalledOnPreviousZone() {
	zone, _ := newZone(":0z", suite.world, suite.logger)
	suite.world.On("Publish", mock.Anything).Return(nil)
	user := &mocks.User{}
	prvZone := &mocks.Zone{}
	prvZone.On("Leave", user).Return(nil, nil)
	user.On("Zone").Return(prvZone)
	user.On("PubSubJSON").Return(&types.UserPubSubJSON{})
	err := zone.Join(user)
	assert.NoError(suite.T(), err)
	prvZone.AssertCalled(suite.T(), "Leave", user)
}

func (suite *ZoneTestSuite) TestJoin_LeaveIsCalledOnPreviousZone_Error() {
	err := errors.New("err")
	zone, _ := newZone(":0z", suite.world, suite.logger)
	user := &mocks.User{}
	prvZone := &mocks.Zone{}
	prvZone.On("Leave", user).Return(err)
	user.On("Zone").Return(prvZone)
	actualErr := zone.Join(user)
	assert.Equal(suite.T(), err, actualErr)
}

func (suite *ZoneTestSuite) TestJoin_PubSubJoin_Error() {
	zone, _ := newZone(":0z", suite.world, suite.logger)
	suite.world.On("Publish", mock.Anything).Return(nil)
	user := &mocks.User{}
	prvZone := &mocks.Zone{}
	prvZone.On("Leave", user).Return(nil)
	user.On("Zone").Return(prvZone)
	user.On("PubSubJSON").Return(&mocks.PubSubJSON{}) // typecast error
	err := zone.Join(user)
	assert.Error(suite.T(), err)
}

func (suite *ZoneTestSuite) TestJoin_Publish() {
	zone, _ := newZone(":0z", suite.world, suite.logger)
	suite.world.On("Publish", mock.Anything).Return(nil)
	user := &mocks.User{}
	prvZone := &mocks.Zone{}
	prvZone.On("Leave", user).Return(nil)
	user.On("Zone").Return(prvZone)
	user.On("PubSubJSON").Return(&types.UserPubSubJSON{})
	err := zone.Join(user)
	assert.NoError(suite.T(), err)
	suite.world.AssertCalled(suite.T(), "Publish", mock.Anything)
}

func (suite *ZoneTestSuite) TestJoin_Publish_Error() {
	err := errors.New("err")
	zone, _ := newZone(":0z", suite.world, suite.logger)
	suite.world.On("Publish", mock.Anything).Return(err)
	user := &mocks.User{}
	prvZone := &mocks.Zone{}
	prvZone.On("Leave", user).Return(nil)
	user.On("Zone").Return(prvZone)
	user.On("PubSubJSON").Return(&types.UserPubSubJSON{})
	actualErr := zone.Join(user)
	assert.Equal(suite.T(), err, actualErr)
	suite.world.AssertCalled(suite.T(), "Publish", mock.Anything)
}

func (suite *ZoneTestSuite) TestLeave() {
	zone, _ := newZone(":0z", suite.world, suite.logger)
	suite.world.On("Publish", mock.Anything).Return(nil)

	user := &mocks.User{}
	user.On("ID").Return("user1")
	user.On("Zone").Return(zone)

	err := zone.Leave(user)
	assert.NoError(suite.T(), err)
	suite.world.AssertCalled(suite.T(), "Publish", mock.Anything)
}

func (suite *ZoneTestSuite) TestMessage() {
	zone, _ := newZone(":0z", suite.world, suite.logger)
	suite.world.On("Publish", mock.Anything).Return(nil)

	user := &mocks.User{}
	user.On("ID").Return("user1")

	err := zone.Message(user, "test")
	assert.NoError(suite.T(), err)
	suite.world.AssertCalled(suite.T(), "Publish", mock.Anything)
}

func (suite *ZoneTestSuite) TestSplit() {
	zone, _ := newZone(":0z", suite.world, suite.logger)

	user1 := &mocks.User{}
	user1.On("ID").Return("user1")
	user2 := &mocks.User{}
	user2.On("ID").Return("user2")

	suite.world.On("Zones").Return(suite.zones)
	suite.world.On("Users").Return(suite.users)
	suite.zones.On("Save", mock.Anything).Return(nil)
	suite.users.On("User", "user1").Return(user1, nil)
	suite.users.On("User", "user2").Return(user2, nil)

	// omg

	leftZone := &mocks.Zone{}
	rightZone := &mocks.Zone{}
	suite.world.On("FindOpenZone", zone, user1).Return(leftZone, nil)
	suite.world.On("FindOpenZone", zone, user2).Return(rightZone, nil)
	user1.On("SetZone", leftZone).Return(nil)
	user1.On("PubSubJSON").Return(&types.UserPubSubJSON{})
	user2.On("SetZone", rightZone).Return(nil)
	user2.On("PubSubJSON").Return(&types.UserPubSubJSON{})

	leftZone.On("ID").Return(":0g")
	leftZone.On("AddUser", user1).Return(nil)
	leftZone.On("PubSubJSON").Return(&types.ZonePubSubJSON{})
	rightZone.On("ID").Return("hz")
	rightZone.On("AddUser", user2).Return(nil)
	rightZone.On("PubSubJSON").Return(&types.ZonePubSubJSON{})
	suite.db.On("SaveUsersAndZones", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	zone.AddUser(user1)
	zone.AddUser(user2)

	zone.Split()

	assert.False(suite.T(), zone.IsOpen())
	assert.Equal(suite.T(), 0, zone.Count())
	leftZone.AssertCalled(suite.T(), "AddUser", user1)
	user1.AssertCalled(suite.T(), "SetZone", leftZone)
	rightZone.AssertCalled(suite.T(), "AddUser", user2)
	user2.AssertCalled(suite.T(), "SetZone", rightZone)
	suite.db.AssertCalled(suite.T(), "SaveUsersAndZones", mock.Anything, mock.Anything, mock.Anything)
}

func TestZoneTestSuit(t *testing.T) {
	suite.Run(t, new(ZoneTestSuite))
}
