package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/cache"
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var seattle = &LatLng{47.6235616, -122.330341, "c23nb"}
var rome = &LatLng{41.9, 12.5, "sr2yk"}
var invalidZoneID = ":::::::"

type WorldTestSuite struct {
	suite.Suite
	chat   *mocks.Chat
	cache  *mocks.Cache
	zone   *mocks.Zone
	pubsub *mocks.PubSub
	user   *mocks.User
}

func (suite *WorldTestSuite) SetupTest() {
	suite.chat = &mocks.Chat{}
	suite.cache = &mocks.Cache{}
	suite.zone = &mocks.Zone{}
	suite.pubsub = &mocks.PubSub{}
	suite.user = &mocks.User{}
}

func (suite *WorldTestSuite) TestNewWorld() {
	ch := make(<-chan types.Event)
	suite.chat.On("PubSub").Return(suite.pubsub)
	suite.chat.On("Cache").Return(suite.cache)
	suite.cache.On("Zone", rootZoneID).Return(suite.zone, nil)
	suite.pubsub.On("Subscribe").Return(ch)

	world, err := newWorld("worldid", suite.chat, 1)
	defer world.close()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "worldid", world.id)
	assert.Equal(suite.T(), suite.zone, world.root)
	assert.Equal(suite.T(), suite.chat, world.chat)
	assert.Equal(suite.T(), 1, world.maxUsersPerZone)
	assert.Equal(suite.T(), ch, world.subscribe)
}

func (suite *WorldTestSuite) TestNewWorldReturnsError() {
	worldErr := errors.New("err")
	suite.chat.On("PubSub").Return(suite.pubsub)
	suite.chat.On("Cache").Return(suite.cache)
	suite.cache.On("Zone", rootZoneID).Return(nil, worldErr)
	suite.pubsub.On("Subscribe").Return(make(<-chan types.Event))

	world, err := newWorld("", suite.chat, 1)
	defer world.close()
	assert.Nil(suite.T(), world)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), worldErr, err)
}

func (suite *WorldTestSuite) TestGetOrCreateZone() {
	world := &World{chat: suite.chat}
	suite.chat.On("Cache").Return(suite.cache)
	suite.cache.On("Zone", rootZoneID).Return(suite.zone, nil)

	z, err := world.GetOrCreateZone(rootZoneID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.zone, z)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMiss() {
	world := &World{chat: suite.chat}
	suite.chat.On("Cache").Return(suite.cache)
	suite.cache.On("Zone", rootZoneID).Return(nil, nil)
	suite.cache.On("SetZone", mock.Anything).Return(nil)

	z, err := world.GetOrCreateZone(rootZoneID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), rootZoneID, z.ID())
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMissAndNewZoneError() {
	world := &World{chat: suite.chat}
	suite.chat.On("Cache").Return(suite.cache)
	suite.cache.On("Zone", invalidZoneID).Return(nil, nil)

	z, err := world.GetOrCreateZone(invalidZoneID)
	assert.Nil(suite.T(), z)
	assert.Equal(suite.T(), "Invalid id", err.Error())
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMissAndSetZoneError() {
	err := errors.New("err")
	world := &World{chat: suite.chat}
	suite.chat.On("Cache").Return(suite.cache)
	suite.cache.On("Zone", rootZoneID).Return(nil, nil)
	suite.cache.On("SetZone", mock.Anything).Return(err)
	z, zerr := world.GetOrCreateZone(rootZoneID)
	assert.Equal(suite.T(), err, zerr)
	assert.Nil(suite.T(), z)
}

func (suite *WorldTestSuite) TestMultipleWorldsWithSameDBDependencyReturnsSameRoot() {
	suite.chat.On("Cache").Return(suite.cache)
	suite.chat.On("PubSub").Return(suite.pubsub)
	suite.cache.On("Zone", rootZoneID).Return(suite.zone, nil)
	suite.pubsub.On("Subscribe").Return(make(<-chan types.Event))

	world1, err1 := newWorld("", suite.chat, 1)
	defer world1.close()
	world2, err2 := newWorld("", suite.chat, 1)
	defer world2.close()
	world3, err3 := newWorld("", suite.chat, 1)
	defer world3.close()

	assert.Equal(suite.T(), suite.zone, world1.root)
	assert.Equal(suite.T(), suite.zone, world2.root)
	assert.Equal(suite.T(), suite.zone, world3.root)

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NoError(suite.T(), err3)
}

func (suite *WorldTestSuite) TestSetZone() {
	suite.cache.On("SetZone", suite.zone).Return(nil)
	suite.chat.On("Cache").Return(suite.cache)
	world := &World{chat: suite.chat}
	err := world.SetZone(suite.zone)
	assert.NoError(suite.T(), err)
	suite.cache.AssertCalled(suite.T(), "SetZone", suite.zone)
}

func (suite *WorldTestSuite) TestSetZoneReturnsError() {
	err := errors.New("sdfsd")
	suite.cache.On("SetZone", suite.zone).Return(err)
	suite.chat.On("Cache").Return(suite.cache)
	world := &World{chat: suite.chat}
	zerr := world.SetZone(suite.zone)
	assert.Error(suite.T(), zerr)
	suite.cache.AssertCalled(suite.T(), "SetZone", suite.zone)
	assert.Equal(suite.T(), err, zerr)
}

func (suite *WorldTestSuite) TestSetUser() {
	suite.cache.On("SetUser", suite.user).Return(nil)
	suite.chat.On("Cache").Return(suite.cache)
	world := &World{chat: suite.chat}
	err := world.SetUser(suite.user)
	assert.NoError(suite.T(), err)
	suite.cache.AssertCalled(suite.T(), "SetUser", suite.user)
}

func (suite *WorldTestSuite) TestSetUserReturnsError() {
	err := errors.New("sdfsd")
	suite.cache.On("SetUser", suite.user).Return(err)
	suite.chat.On("Cache").Return(suite.cache)
	world := &World{chat: suite.chat}
	zerr := world.SetUser(suite.user)
	assert.Error(suite.T(), zerr)
	suite.cache.AssertCalled(suite.T(), "SetUser", suite.user)
	assert.Equal(suite.T(), err, zerr)
}

func TestWorldSuite(t *testing.T) {
	suite.Run(t, new(WorldTestSuite))
}

type CreateZoneForUserSuite struct {
	suite.Suite
	chat   *mocks.Chat
	cache  *mocks.Cache
	root   *mocks.Zone
	pubsub *mocks.PubSub
	user   *mocks.User
	left   *mocks.Zone
	right  *mocks.Zone
}

func (suite *CreateZoneForUserSuite) SetupTest() {
	suite.chat = &mocks.Chat{}
	suite.cache = &mocks.Cache{}
	suite.root = &mocks.Zone{}
	suite.pubsub = &mocks.PubSub{}
	suite.user = &mocks.User{}
	suite.left = &mocks.Zone{}
	suite.right = &mocks.Zone{}

	suite.chat.On("Cache").Return(suite.cache)
	suite.cache.On("Zone", rootZoneID).Return(suite.root, nil)
	suite.chat.On("PubSub").Return(suite.pubsub)
	suite.root.On("Geohash").Return("")
	suite.root.On("LeftZoneID").Return(":0g")
	suite.root.On("RightZoneID").Return(":hz")
	suite.pubsub.On("Subscribe").Return(make(<-chan types.Event))
	suite.right.On("Geohash").Return("")
	suite.right.On("From").Return("h")
}

func (suite *CreateZoneForUserSuite) TestGetOrCreateZoneForUser_EmptyWorld() {
	suite.root.On("IsOpen").Return(true)
	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	zone, err := world.GetOrCreateZoneForUser(suite.user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), world.root, zone)
}

func (suite *CreateZoneForUserSuite) TestGetOrCreateZoneForUser_RightZoneReturnsError() {
	err := errors.New("invalid")
	suite.root.On("IsOpen").Return(false)
	suite.right.On("IsOpen").Return(true)
	suite.user.On("Location").Return(rome)
	suite.cache.On("Zone", ":hz").Return(suite.right, err)

	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	zone, zerr := world.GetOrCreateZoneForUser(suite.user)
	assert.Equal(suite.T(), err, zerr)
	assert.Nil(suite.T(), zone)
}

func (suite *CreateZoneForUserSuite) TestGetOrCreateZoneForUser_LeftZoneReturnsError() {
	err := errors.New("invalid")
	suite.root.On("IsOpen").Return(false)
	suite.left.On("IsOpen").Return(true)
	suite.user.On("Location").Return(seattle)
	suite.cache.On("Zone", ":hz").Return(suite.right, nil)
	suite.cache.On("Zone", ":0g").Return(suite.left, err)

	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	zone, zerr := world.GetOrCreateZoneForUser(suite.user)
	assert.Equal(suite.T(), err, zerr)
	assert.Nil(suite.T(), zone)
}

func (suite *CreateZoneForUserSuite) TestGetOrCreateZoneForUser_ReturnsLeftZone() {
	suite.root.On("IsOpen").Return(false)
	suite.left.On("IsOpen").Return(true)
	suite.user.On("Location").Return(seattle)
	suite.cache.On("Zone", ":hz").Return(suite.right, nil)
	suite.cache.On("Zone", ":0g").Return(suite.left, nil)

	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	zone, zerr := world.GetOrCreateZoneForUser(suite.user)
	assert.NoError(suite.T(), zerr)
	assert.Equal(suite.T(), suite.left, zone)
}

func (suite *CreateZoneForUserSuite) TestGetOrCreateZoneForUser_ReturnsRightZone() {
	suite.root.On("IsOpen").Return(false)
	suite.right.On("IsOpen").Return(true)
	suite.user.On("Location").Return(rome)
	suite.cache.On("Zone", ":hz").Return(suite.right, nil)
	suite.cache.On("Zone", ":0g").Return(suite.left, nil)

	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	zone, zerr := world.GetOrCreateZoneForUser(suite.user)
	assert.NoError(suite.T(), zerr)
	assert.Equal(suite.T(), suite.right, zone)
}

func (suite *CreateZoneForUserSuite) TestGetOrCreateZoneForUser_ReturnsLeftZone2() {
	suite.root.On("IsOpen").Return(false)
	suite.left.On("IsOpen").Return(true)
	suite.user.On("Location").Return(seattle)
	suite.cache.On("Zone", ":0g").Return(suite.left, nil)

	right := &mocks.Zone{}
	right.On("Geohash").Return("s")
	right.On("From").Return("h")
	suite.cache.On("Zone", ":hz").Return(right, nil)

	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	zone, zerr := world.GetOrCreateZoneForUser(suite.user)
	assert.NoError(suite.T(), zerr)
	assert.Equal(suite.T(), suite.left, zone)
}

func (suite *CreateZoneForUserSuite) TestGetOrCreateZoneForUser_ReturnsRightZone2() {
	suite.root.On("IsOpen").Return(false)
	suite.user.On("Location").Return(rome)
	suite.cache.On("Zone", ":0g").Return(suite.left, nil)

	right := &mocks.Zone{}
	right.On("Geohash").Return("s")
	right.On("From").Return("h")
	right.On("IsOpen").Return(true)
	suite.cache.On("Zone", ":hz").Return(right, nil)

	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	zone, zerr := world.GetOrCreateZoneForUser(suite.user)
	assert.NoError(suite.T(), zerr)
	assert.Equal(suite.T(), right, zone)
}

func (suite *CreateZoneForUserSuite) TestGetOrCreateZoneForUser_ErrorOnNoOpenRooms() {
	suite.root.On("IsOpen").Return(false)
	suite.cache.On("Zone", ":0g").Return(suite.left, nil)
	suite.cache.On("Zone", ":hz").Return(suite.right, nil)
	suite.left.On("IsOpen").Return(false)
	suite.right.On("IsOpen").Return(false)
	latlng := &mocks.LatLng{}
	latlng.On("Geohash").Return("")
	suite.user.On("Location").Return(latlng)

	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	zone, err := world.GetOrCreateZoneForUser(suite.user)
	assert.Nil(suite.T(), zone)
	assert.Equal(suite.T(), "Unable to find zone", err.Error())
}

func TestCreateZoneForUserSuite(t *testing.T) {
	suite.Run(t, new(CreateZoneForUserSuite))
}

type WorldIntegrationTestSuite struct {
	suite.Suite
	world  *World
	db     *mocks.DB
	chat   *mocks.Chat
	cache  *cache.Cache
	pubsub *mocks.PubSub
	ch     chan types.Event
	rch    <-chan types.Event
}

func (suite *WorldIntegrationTestSuite) SetupTest() {
	suite.db = &mocks.DB{}
	suite.chat = &mocks.Chat{}
	suite.cache = cache.NewCache(suite.db)
	suite.pubsub = &mocks.PubSub{}
	suite.ch = make(chan types.Event)
	suite.rch = suite.ch

	suite.db.On("GetZone", mock.Anything).Return(nil, nil)
	suite.db.On("SetZone", mock.Anything).Return(nil)
	suite.chat.On("Cache").Return(suite.cache)
	suite.chat.On("PubSub").Return(suite.pubsub)
	suite.pubsub.On("Subscribe").Return(suite.rch)
}

func (suite *WorldIntegrationTestSuite) TestIntegration() {
	testCases := []string{"000", "z0z", "2k1", "bbc", "zzz", "c23nb"}

	for _, test := range testCases {
		world, err := newWorld("", suite.chat, 1)
		defer world.close()
		assert.NoError(suite.T(), err)
		world.root.SetIsOpen(false)

		user := &mocks.User{}
		user.On("Location").Return(&LatLng{0, 0, test})

		var zone types.Zone
		for i := 0; i < 5*len(test); i++ {
			zone, err = world.GetOrCreateZoneForUser(user)
			assert.NoError(suite.T(), err)
			zone.SetIsOpen(false)
			if len(zone.Geohash()) == len(test) {
				break
			}
		}
		assert.Equal(suite.T(), test+rootZoneID, zone.ID())
	}
}

// This is really gross :-(
func (suite *WorldIntegrationTestSuite) TestIncomingEventsCallOnReceive() {

	world, err := newWorld("", suite.chat, 1)
	defer world.close()
	assert.NoError(suite.T(), err)

	mockEvent := &mocks.Event{}
	mockEventData := &mocks.EventData{}

	done := make(chan bool)
	mockEventData.On("OnReceive", mockEvent).Return(nil).Run(func(args mock.Arguments) {
		assert.Equal(suite.T(), mockEvent, args.Get(0))
		done <- true
	})

	mockEvent.On("Data").Return(mockEventData)
	mockEvent.On("SetWorld", world).Return()
	mockEventData.AssertNotCalled(suite.T(), "OnReceive", mockEvent)
	suite.ch <- mockEvent
	<-done
	suite.pubsub.AssertCalled(suite.T(), "Subscribe")
	mockEventData.AssertCalled(suite.T(), "OnReceive", mockEvent)
}

func TestWorldIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(WorldIntegrationTestSuite))
}

type PubSubSuite struct {
	suite.Suite
	chat   *mocks.Chat
	cache  *mocks.Cache
	pubsub *mocks.PubSub
	events *mocks.EventFactory
	event  *mocks.Event
	data   *mocks.EventData
	root   *mocks.Zone
}

func (suite *PubSubSuite) SetupTest() {
	suite.chat = &mocks.Chat{}
	suite.cache = &mocks.Cache{}
	suite.pubsub = &mocks.PubSub{}
	suite.event = &mocks.Event{}
	suite.events = &mocks.EventFactory{}
	suite.data = &mocks.EventData{}
	suite.root = &mocks.Zone{}

	suite.chat.On("Cache").Return(suite.cache)
	suite.chat.On("PubSub").Return(suite.pubsub)
	suite.chat.On("Events").Return(suite.events)
	suite.cache.On("Zone", rootZoneID).Return(suite.root, nil)
	suite.event.On("Data").Return(suite.data)
	suite.pubsub.On("Subscribe").Return(make(<-chan types.Event))
	suite.events.On("New", "", suite.data).Return(suite.event, nil)
}

func (suite *PubSubSuite) TestPublishCallsBeforePublish() {
	suite.pubsub.On("Publish", suite.event).Return(nil)
	suite.data.On("BeforePublish", suite.event).Return(nil)

	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	err := world.Publish(suite.data)
	assert.NoError(suite.T(), err)
	suite.data.AssertCalled(suite.T(), "BeforePublish", suite.event)
}

func (suite *PubSubSuite) TestPublishReturnsBeforePublishError() {
	err1 := errors.New("err")
	suite.pubsub.On("Publish", suite.event).Return(nil)
	suite.data.On("BeforePublish", suite.event).Return(err1)

	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	err2 := world.Publish(suite.data)
	assert.Equal(suite.T(), err1, err2)
	suite.data.AssertCalled(suite.T(), "BeforePublish", suite.event)
}

func (suite *PubSubSuite) TestPublishReturnsPubSubError() {
	err1 := errors.New("err")
	suite.pubsub.On("Publish", suite.event).Return(err1)
	suite.data.On("BeforePublish", suite.event).Return(nil)

	world, _ := newWorld("", suite.chat, 1)
	defer world.close()
	err2 := world.Publish(suite.data)
	assert.Equal(suite.T(), err1, err2)
	suite.data.AssertCalled(suite.T(), "BeforePublish", suite.event)
}

func TestPubSubSuite(t *testing.T) {
	suite.Run(t, new(PubSubSuite))
}
