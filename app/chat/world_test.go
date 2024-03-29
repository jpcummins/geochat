package chat

import (
	// "errors"
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var seattle = newLatLng(47.6235616, -122.330341)
var rome = newLatLng(41.9, 12.5)
var empty = newLatLng(0, 0)

var invalidZoneID = ":::::::"

var leftZoneID = ":0g"
var rightZoneID = ":hz"

type WorldTestSuite struct {
	suite.Suite
	world  *World
	root   *mocks.Zone
	zones  *mocks.Zones
	user   *mocks.User
	users  *mocks.Users
	db     *mocks.DB
	pubsub *mocks.PubSub
	left   *mocks.Zone
	right  *mocks.Zone
	logger *mocks.Logger
	ch     chan types.PubSubEvent
	rch    <-chan types.PubSubEvent
}

func (suite *WorldTestSuite) SetupTest() {
	suite.ch = make(chan types.PubSubEvent)
	suite.rch = suite.ch
	suite.zones = &mocks.Zones{}
	suite.user = &mocks.User{}
	suite.users = &mocks.Users{}
	suite.db = &mocks.DB{}
	suite.pubsub = &mocks.PubSub{}
	suite.root = &mocks.Zone{}
	suite.left = &mocks.Zone{}
	suite.right = &mocks.Zone{}
	suite.logger = &mocks.Logger{}
	suite.world = suite.NewWorld()
}

func (suite *WorldTestSuite) NewWorld() *World {
	world := &World{
		WorldPubSubJSON: &types.WorldPubSubJSON{
			ID:       rootWorldID,
			MaxUsers: 10,
		},
		db:     suite.db,
		pubsub: suite.pubsub,
		users:  suite.users,
		zones:  suite.zones,
		root:   suite.root,
		logger: suite.logger,
	}
	suite.root.On("Geohash").Return("")
	suite.root.On("RightZoneID").Return(rightZoneID)
	suite.root.On("LeftZoneID").Return(leftZoneID)
	suite.right.On("Geohash").Return("")
	suite.right.On("From").Return("h")
	suite.pubsub.On("Subscribe").Return(suite.rch)
	suite.logger.On("New", mock.Anything).Return(suite.logger)
	return world
}

func (suite *WorldTestSuite) TestNewWorld() {
	suite.db.On("Zone", rootZoneID, mock.Anything, mock.Anything).Return(nil, nil)
	suite.db.On("SaveZone", mock.Anything, mock.Anything).Return(nil)

	world, err := newWorld(rootWorldID, suite.db, suite.pubsub, suite.logger)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), rootWorldID, world.ID())
	assert.Equal(suite.T(), rootZoneID, world.root.ID())
	assert.Equal(suite.T(), suite.db, world.db)
	assert.Equal(suite.T(), suite.pubsub, world.pubsub)
}

//
// func (suite *WorldTestSuite) TestInitReturnsError() {
// 	worldErr := errors.New("err")
// 	suite.db.On("Zone", rootZoneID, mock.Anything, mock.Anything).Return(nil, worldErr)
//
// 	_, err := newWorld(rootWorldID, suite.db, suite.pubsub, 10)
// 	assert.Equal(suite.T(), worldErr, err)
// }
//
// func (suite *WorldTestSuite) TestGetOrCreateZone() {
// 	suite.zones.On("Zone", rootZoneID).Return(suite.root, nil)
// 	suite.root.On("ID").Return(rootZoneID)
//
// 	z, err := suite.world.GetOrCreateZone(rootZoneID)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), rootZoneID, z.ID())
// }
//
// func (suite *WorldTestSuite) TestGetOrCreateZoneError() {
// 	mockError := errors.New("invalid id")
// 	suite.root.On("ID").Return(rootZoneID)
// 	suite.zones.On("Zone", invalidZoneID).Return(nil, mockError)
//
// 	z, err := suite.world.GetOrCreateZone(invalidZoneID)
// 	assert.Nil(suite.T(), z)
// 	assert.Equal(suite.T(), mockError, err)
// }
//
// func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMissAndSaveZoneError() {
// 	err := errors.New("err")
// 	suite.zones.On("Zone", rootZoneID).Return(nil, nil)
// 	suite.zones.On("Save", mock.Anything).Return(err)
//
// 	z, zerr := suite.world.GetOrCreateZone(rootZoneID)
// 	assert.Equal(suite.T(), err, zerr)
// 	assert.Nil(suite.T(), z)
// }
//
// func (suite *WorldTestSuite) TestMultipleWorldsWithSameDBDependencyReturnsSameRoot() {
// 	suite.zones.On("Zone", rootZoneID).Return(suite.root, nil)
// 	suite.root.On("ID").Return(rootZoneID)
//
// 	zone1, _ := suite.NewWorld().GetOrCreateZone(rootZoneID)
// 	zone2, _ := suite.NewWorld().GetOrCreateZone(rootZoneID)
// 	zone3, _ := suite.NewWorld().GetOrCreateZone(rootZoneID)
//
// 	assert.Equal(suite.T(), suite.root, zone1)
// 	assert.Equal(suite.T(), suite.root, zone2)
// 	assert.Equal(suite.T(), suite.root, zone3)
// }
//
// func (suite *WorldTestSuite) TestFindOpenZone_EmptyWorld() {
// 	suite.root.On("IsOpen").Return(true)
//
// 	zone, err := suite.world.FindOpenZone(suite.user)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), suite.root, zone)
// }
//
// func (suite *WorldTestSuite) TestFindOpenZone_RightZoneReturnsError() {
// 	err := errors.New("invalid")
// 	suite.root.On("IsOpen").Return(false)
// 	suite.user.On("Location").Return(seattle)
// 	suite.zones.On("Zone", rightZoneID).Return(nil, err)
//
// 	zone, zerr := suite.world.FindOpenZone(suite.user)
// 	assert.Equal(suite.T(), err, zerr)
// 	assert.Nil(suite.T(), zone)
// }
//
// func (suite *WorldTestSuite) TestFindOpenZone_LeftZoneReturnsError() {
// 	err := errors.New("invalid")
// 	suite.root.On("IsOpen").Return(false)
// 	suite.user.On("Location").Return(seattle)
// 	suite.zones.On("Zone", rightZoneID).Return(suite.right, nil)
// 	suite.zones.On("Zone", leftZoneID).Return(nil, err)
//
// 	zone, zerr := suite.world.FindOpenZone(suite.user)
// 	assert.Equal(suite.T(), err, zerr)
// 	assert.Nil(suite.T(), zone)
// }
//
// func (suite *WorldTestSuite) TestFindOpenZone_ReturnsLeftZone() {
// 	suite.root.On("IsOpen").Return(false)
// 	suite.left.On("IsOpen").Return(true)
// 	suite.user.On("Location").Return(seattle)
// 	suite.zones.On("Zone", rightZoneID).Return(suite.right, nil)
// 	suite.zones.On("Zone", leftZoneID).Return(suite.left, nil)
//
// 	zone, zerr := suite.world.FindOpenZone(suite.user)
// 	assert.Nil(suite.T(), zerr)
// 	assert.Equal(suite.T(), suite.left, zone)
// }
//
// func (suite *WorldTestSuite) TestFindOpenZone_ReturnsRightZone() {
// 	suite.root.On("IsOpen").Return(false)
// 	suite.right.On("IsOpen").Return(true)
// 	suite.user.On("Location").Return(rome)
// 	suite.zones.On("Zone", rightZoneID).Return(suite.right, nil)
// 	suite.zones.On("Zone", leftZoneID).Return(suite.left, nil)
//
// 	zone, zerr := suite.world.FindOpenZone(suite.user)
// 	assert.Nil(suite.T(), zerr)
// 	assert.Equal(suite.T(), suite.right, zone)
// }
//
// func (suite *WorldTestSuite) TestFindOpenZone_ReturnsLeftZone2() {
// 	suite.root.On("IsOpen").Return(false)
// 	suite.left.On("IsOpen").Return(true)
// 	suite.user.On("Location").Return(seattle)
// 	suite.zones.On("Zone", rightZoneID).Return(suite.right, nil)
// 	suite.zones.On("Zone", leftZoneID).Return(suite.left, nil)
// 	suite.right.On("Geohash").Return("s")
// 	suite.right.On("From").Return("h")
//
// 	zone, zerr := suite.world.FindOpenZone(suite.user)
// 	assert.NoError(suite.T(), zerr)
// 	assert.Equal(suite.T(), suite.left, zone)
// }
//
// func (suite *WorldTestSuite) TestFindOpenZone_ReturnsRightZone2() {
// 	suite.root.On("IsOpen").Return(false)
// 	suite.right.On("IsOpen").Return(true)
// 	suite.user.On("Location").Return(rome)
// 	suite.zones.On("Zone", rightZoneID).Return(suite.right, nil)
// 	suite.zones.On("Zone", leftZoneID).Return(suite.left, nil)
// 	suite.right.On("Geohash").Return("s")
// 	suite.right.On("From").Return("h")
//
// 	zone, zerr := suite.world.FindOpenZone(suite.user)
// 	assert.NoError(suite.T(), zerr)
// 	assert.Equal(suite.T(), suite.right, zone)
// }
//
// func (suite *WorldTestSuite) TestFindOpenZone_ErrorOnNoOpenRooms() {
// 	suite.root.On("IsOpen").Return(false)
// 	suite.right.On("IsOpen").Return(false)
// 	suite.user.On("Location").Return(empty)
// 	suite.zones.On("Zone", rightZoneID).Return(suite.right, nil)
// 	suite.zones.On("Zone", leftZoneID).Return(suite.left, nil)
// 	suite.right.On("Geohash").Return("s")
// 	suite.right.On("From").Return("h")
//
// 	zone, err := suite.world.FindOpenZone(suite.user)
// 	assert.Nil(suite.T(), zone)
// 	assert.Equal(suite.T(), "Unable to find zone", err.Error())
// }
//
// func (suite *WorldTestSuite) TestIntegration() {
// 	suite.db.On("Zone", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
// 	suite.db.On("SaveZone", mock.Anything, mock.Anything).Return(nil)
// 	suite.pubsub.On("Subscribe").Return(nil)
//
// 	testCases := []string{"000", "z0z", "2k1", "bbc", "zzz", "c23nb"}
//
// 	for _, test := range testCases {
// 		world, err := newWorld(rootWorldID, suite.db, suite.pubsub, 1, suite.logger)
// 		assert.NoError(suite.T(), err)
// 		world.root.SetIsOpen(false)
//
// 		user := &mocks.User{}
// 		user.On("Location").Return(&LatLng{geohash: test})
//
// 		var zone types.Zone
// 		for i := 0; i < 5*len(test); i++ {
// 			zone, err = world.FindOpenZone(world.Zone(), user)
// 			assert.NoError(suite.T(), err)
// 			zone.SetIsOpen(false)
// 			if len(zone.Geohash()) == len(test) {
// 				break
// 			}
// 		}
// 		assert.Equal(suite.T(), test+rootZoneID, zone.ID())
// 	}
// }

//
// func (suite *WorldTestSuite) TestIncomingEventsCallOnReceive() {
// 	suite.db.On("Zone", rootZoneID, mock.Anything, mock.Anything).Return(nil, nil)
// 	suite.db.On("SaveZone", mock.Anything, mock.Anything).Return(nil)
//
// 	mockEvent := &mocks.PubSubEvent{}
// 	mockEventData := &mocks.PubSubEventData{}
//
// 	done := make(chan bool)
// 	mockEventData.On("OnReceive", mockEvent).Return(nil).Run(func(args mock.Arguments) {
// 		assert.Equal(suite.T(), mockEvent, args.Get(0))
// 		done <- true
// 	})
//
// 	mockEvent.On("Data").Return(mockEventData)
// 	mockEvent.On("SetWorld", mock.Anything).Return()
//
// 	newWorld(rootWorldID, suite.db, suite.pubsub, 10)
//
// 	mockEventData.AssertNotCalled(suite.T(), "OnReceive", mockEvent)
// 	suite.ch <- mockEvent
// 	<-done
// 	suite.pubsub.AssertCalled(suite.T(), "Subscribe")
// 	mockEventData.AssertCalled(suite.T(), "OnReceive", mockEvent)
// }
//
// func (suite *WorldTestSuite) TestPublishCallsBeforePublish() {
// 	event := &mocks.PubSubEvent{}
// 	data := &mocks.PubSubEventData{}
//
// 	event.On("Data").Return(data)
// 	suite.pubsub.On("Publish", event).Return(nil)
// 	data.On("BeforePublish", event).Return(nil)
//
// 	err := suite.world.Publish(event)
// 	assert.NoError(suite.T(), err)
// 	data.AssertCalled(suite.T(), "BeforePublish", event)
// }
//
// func (suite *WorldTestSuite) TestPublishReturnsBeforePublishError() {
// 	err1 := errors.New("err")
// 	event := &mocks.PubSubEvent{}
// 	data := &mocks.PubSubEventData{}
// 	event.On("Data").Return(data)
//
// 	data.On("BeforePublish", event).Return(err1)
// 	err2 := suite.world.Publish(event)
// 	assert.Equal(suite.T(), err1, err2)
// }
//
// func (suite *WorldTestSuite) TestPublishReturnsPubSubError() {
// 	err1 := errors.New("err")
// 	event := &mocks.PubSubEvent{}
// 	data := &mocks.PubSubEventData{}
// 	event.On("Data").Return(data)
// 	suite.pubsub.On("Publish", event).Return(err1)
// 	data.On("BeforePublish", event).Return(nil)
// 	err2 := suite.world.Publish(event)
// 	assert.Equal(suite.T(), err1, err2)
// }
//
// func (suite *WorldTestSuite) TestNewUser() {
// 	suite.users.On("Save", mock.Anything).Return(nil)
//
// 	world := suite.NewWorld()
// 	user, err := world.NewUser("abc", "bob", seattle.lat, seattle.lng)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), "abc", user.ID())
// 	assert.Equal(suite.T(), "bob", user.Name())
// 	assert.Equal(suite.T(), seattle.lat, user.Location().Lat())
// 	assert.Equal(suite.T(), seattle.lng, user.Location().Lng())
// 	assert.Nil(suite.T(), user.Zone())
// }
//
// func (suite *WorldTestSuite) TestNewUserError() {
// 	err1 := errors.New("err")
// 	suite.users.On("Save", mock.Anything).Return(err1)
//
// 	world := suite.NewWorld()
// 	user, err2 := world.NewUser("abc", "bob", seattle.lat, seattle.lng)
// 	assert.Equal(suite.T(), err1, err2)
// 	assert.Nil(suite.T(), user)
// }
//

func (suite *WorldTestSuite) TestParentID_Root() {
	suite.db.On("Zone", rootZoneID, mock.Anything, mock.Anything).Return(nil, nil)
	suite.db.On("SaveZone", mock.Anything, mock.Anything).Return(nil)
	world, _ := newWorld(rootWorldID, suite.db, suite.pubsub, suite.logger)
	assert.Equal(suite.T(), "", world.Zone().ParentZoneID())
}

func TestWorldSuite(t *testing.T) {
	suite.Run(t, new(WorldTestSuite))
}
