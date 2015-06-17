package chat

import (
	"errors"
	// "github.com/jpcummins/geochat/app/cache"
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	// "testing"
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
	suite.cache.On("Zone", ":0z").Return(suite.zone)
	suite.pubsub.On("Subscribe").Return(ch)
	world, err := newWorld("worldid", suite.chat, 1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "worldid", world.id)
	assert.Equal(suite.T(), suite.zone, world.root)
	assert.Equal(suite.T(), suite.chat, world.chat)
	assert.Equal(suite.T(), 1, world.maxUsersPerZone)
	assert.Equal(suite.T(), ch, world.subscribe)
}

func (suite *WorldTestSuite) TestNewWorldReturnsError() {
	worldErr := errors.New("err")
	suite.cache.On("Zone", ":0z").Return(nil, worldErr)
	world, err := newWorld("", suite.chat, 1)
	assert.Nil(suite.T(), world)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), worldErr, err)
}

func (suite *WorldTestSuite) TestGetOrCreateZone() {
	suite.cache.On("Zone", ":0z").Return(suite.zone, nil)
	z, err := (&World{}).GetOrCreateZone(":0z")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.zone, z)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMiss() {
	suite.cache.On("Zone", ":0z").Return(nil, nil)
	suite.cache.On("SetZone", mock.Anything).Return(nil)
	z, err := (&World{}).GetOrCreateZone(":0z")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), ":0z", z.ID())
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMissAndNewZoneError() {
	suite.cache.On("Zone", invalidZoneID).Return(nil, nil)
	z, err := (&World{}).GetOrCreateZone(invalidZoneID)
	assert.Nil(suite.T(), z)
	assert.Equal(suite.T(), "Invalid id", err.Error())
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMissAndSetZoneError() {
	err := errors.New("err")
	suite.cache.On("Zone", ":0z").Return(nil, nil)
	suite.cache.On("SetZone", mock.Anything).Return(err)
	z, zerr := (&World{}).GetOrCreateZone(":0z")
	assert.Equal(suite.T(), err, zerr)
	assert.Nil(suite.T(), z)
}

func (suite *WorldTestSuite) TestMultipleWorldsWithSameDBDependencyReturnsSameRoot() {
	suite.cache.On("Zone", ":0z").Return(suite.zone, nil)

	world1, err1 := newWorld("", suite.chat, 1)
	world2, err2 := newWorld("", suite.chat, 1)
	world3, err3 := newWorld("", suite.chat, 1)

	assert.Equal(suite.T(), suite.zone, world1.root)
	assert.Equal(suite.T(), suite.zone, world2.root)
	assert.Equal(suite.T(), suite.zone, world3.root)

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NoError(suite.T(), err3)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_EmptyWorld() {
	suite.cache.On("Zone", ":0z").Return(suite.zone, nil)
	world := &World{}
	zone, err := world.GetOrCreateZoneForUser(&User{})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), world.root, zone)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_RightZoneReturnsError() {
	suite.cache.On("Zone", ":0z").Return(suite.zone, nil)
	suite.zone.On("IsOpen").Return(false)
	suite.user.On("Location").Return(seattle)
	suite.zone.On("RightZoneID").Return(invalidZoneID)
	world, err := newWorld("", suite.chat, 1)
	zone, err := world.GetOrCreateZoneForUser(suite.user)
	assert.Equal(suite.T(), "Invalid id", err)
	assert.Nil(suite.T(), zone)
}

//
// func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_LeftZoneReturnsError() {
// 	err := errors.New("err")
//
// 	root := &mocks.Zone{}
// 	root.On("IsOpen").Return(false)
// 	root.On("Geohash").Return("")
// 	root.On("RightZoneID").Return(":hz")
// 	root.On("LeftZoneID").Return(":0g")
//
// 	suite.cache.On("Zone", ":hz").Return(&mocks.Zone{}, nil)
// 	suite.cache.On("Zone", ":0g").Return(nil, err)
//
// 	user := &mocks.User{}
// 	user.On("Location").Return(seattle)
//
// 	world := &World{
// 		cache:           suite.cache,
// 		root:            root,
// 		maxUsersPerZone: 2,
// 	}
//
// 	zone, err := world.GetOrCreateZoneForUser(user)
// 	assert.Error(suite.T(), err)
// 	assert.Nil(suite.T(), zone)
// }
//
// func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_ReturnsLeftZone() {
// 	root := &mocks.Zone{}
// 	root.On("IsOpen").Return(false)
// 	root.On("Geohash").Return("")
// 	root.On("RightZoneID").Return(":hz")
// 	root.On("LeftZoneID").Return(":0g")
//
// 	leftZone := &mocks.Zone{}
// 	leftZone.On("IsOpen").Return(true)
//
// 	rightZone := &mocks.Zone{}
// 	rightZone.On("Geohash").Return("")
// 	rightZone.On("From").Return("h")
//
// 	suite.cache.On("Zone", ":0g").Return(leftZone, nil)
// 	suite.cache.On("Zone", ":hz").Return(rightZone, nil)
//
// 	user := &mocks.User{}
// 	user.On("Location").Return(seattle)
//
// 	world := &World{
// 		cache:           suite.cache,
// 		root:            root,
// 		maxUsersPerZone: 2,
// 	}
//
// 	zone, err := world.GetOrCreateZoneForUser(user)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), leftZone, zone)
// }
//
// func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_ReturnsRightZone() {
// 	root := &mocks.Zone{}
// 	root.On("IsOpen").Return(false)
// 	root.On("Geohash").Return("")
// 	root.On("RightZoneID").Return(":hz")
// 	root.On("LeftZoneID").Return(":0g")
//
// 	leftZone := &mocks.Zone{}
//
// 	rightZone := &mocks.Zone{}
// 	rightZone.On("Geohash").Return("")
// 	rightZone.On("From").Return("h")
// 	rightZone.On("IsOpen").Return(true)
//
// 	suite.cache.On("Zone", ":0g").Return(leftZone, nil)
// 	suite.cache.On("Zone", ":hz").Return(rightZone, nil)
//
// 	user := &mocks.User{}
// 	user.On("Location").Return(rome)
//
// 	world := &World{
// 		cache:           suite.cache,
// 		root:            root,
// 		maxUsersPerZone: 2,
// 	}
//
// 	zone, err := world.GetOrCreateZoneForUser(user)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), rightZone, zone)
// }
//
// func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_ReturnsLeftZone2() {
// 	root := &mocks.Zone{}
// 	root.On("IsOpen").Return(false)
// 	root.On("Geohash").Return("")
// 	root.On("RightZoneID").Return(":hz")
// 	root.On("LeftZoneID").Return(":0g")
//
// 	leftZone := &mocks.Zone{}
// 	leftZone.On("IsOpen").Return(true)
//
// 	rightZone := &mocks.Zone{}
// 	rightZone.On("Geohash").Return("s")
// 	rightZone.On("From").Return("h")
//
// 	suite.cache.On("Zone", ":0g").Return(leftZone, nil)
// 	suite.cache.On("Zone", ":hz").Return(rightZone, nil)
//
// 	user := &mocks.User{}
// 	user.On("Location").Return(seattle)
//
// 	world := &World{
// 		cache:           suite.cache,
// 		root:            root,
// 		maxUsersPerZone: 2,
// 	}
//
// 	zone, err := world.GetOrCreateZoneForUser(user)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), leftZone, zone)
// }
//
// func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_ReturnsRightZone2() {
// 	root := &mocks.Zone{}
// 	root.On("IsOpen").Return(false)
// 	root.On("Geohash").Return("")
// 	root.On("RightZoneID").Return(":hz")
// 	root.On("LeftZoneID").Return(":0g")
//
// 	leftZone := &mocks.Zone{}
//
// 	rightZone := &mocks.Zone{}
// 	rightZone.On("Geohash").Return("s")
// 	rightZone.On("From").Return("h")
// 	rightZone.On("IsOpen").Return(true)
//
// 	suite.cache.On("Zone", ":0g").Return(leftZone, nil)
// 	suite.cache.On("Zone", ":hz").Return(rightZone, nil)
//
// 	user := &mocks.User{}
// 	user.On("Location").Return(rome)
//
// 	world := &World{
// 		cache:           suite.cache,
// 		root:            root,
// 		maxUsersPerZone: 2,
// 	}
//
// 	zone, err := world.GetOrCreateZoneForUser(user)
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), rightZone, zone)
// }
//
// func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_ErrorOnNoOpenRooms() {
// 	root := &mocks.Zone{}
// 	root.On("IsOpen").Return(false)
// 	root.On("Geohash").Return(seattle.Geohash())
//
// 	user := &mocks.User{}
// 	user.On("Location").Return(seattle)
//
// 	world := &World{
// 		cache:           suite.cache,
// 		root:            root,
// 		maxUsersPerZone: 2,
// 	}
//
// 	zone, err := world.GetOrCreateZoneForUser(user)
// 	assert.Nil(suite.T(), zone)
// 	assert.Error(suite.T(), err)
// 	assert.Equal(suite.T(), "Unable to find zone", err.Error())
// }
//
// func (suite *WorldTestSuite) TestIntegration() {
// 	db := &mocks.DB{}
// 	db.On("GetZone", mock.Anything).Return(nil, nil)
// 	db.On("SetZone", mock.Anything).Return(nil)
//
// 	world := &World{
// 		id:              "test",
// 		cache:           cache.NewCache(db),
// 		factory:         &Factory{},
// 		maxUsersPerZone: 1,
// 	}
//
// 	w, err := world.GetOrCreateZone(":0z")
// 	assert.NoError(suite.T(), err)
// 	world.root = w
// 	world.root.SetIsOpen(false)
// 	assert.Equal(suite.T(), ":0z", world.root.ID())
//
// 	user := &mocks.User{}
// 	user.On("Location").Return(seattle)
//
// 	zone, err := world.GetOrCreateZoneForUser(user)
// 	assert.Equal(suite.T(), ":0g", zone.ID())
//
// 	zone.SetIsOpen(false)
// 	zone, err = world.GetOrCreateZoneForUser(user)
// 	assert.Equal(suite.T(), ":8g", zone.ID())
// }
//
// func (suite *WorldTestSuite) TestIntegration2() {
// 	testCases := []string{"000", "z0z", "2k1", "bbc", "zzz", "c23nb"}
//
// 	db := &mocks.DB{}
// 	db.On("GetZone", mock.Anything).Return(nil, nil)
// 	db.On("SetZone", mock.Anything).Return(nil)
//
// 	for _, test := range testCases {
// 		world := &World{
// 			cache:           cache.NewCache(db),
// 			factory:         &Factory{},
// 			maxUsersPerZone: 1,
// 		}
// 		w, err := world.GetOrCreateZone(":0z")
// 		assert.NoError(suite.T(), err)
// 		world.root = w
// 		world.root.SetIsOpen(false)
//
// 		user := &mocks.User{}
// 		user.On("Location").Return(&LatLng{0, 0, test})
//
// 		var zone types.Zone
// 		for i := 0; i < 5*len(test); i++ {
// 			zone, err = world.GetOrCreateZoneForUser(user)
// 			assert.NoError(suite.T(), err)
// 			zone.SetIsOpen(false)
// 			if len(zone.Geohash()) == len(test) {
// 				break
// 			}
// 		}
// 		assert.Equal(suite.T(), test+":0z", zone.ID())
// 	}
// }
//
// func (suite *WorldTestSuite) TestPublishCallsBeforePublish() {
// 	event := &mocks.Event{}
// 	data := &mocks.EventData{}
// 	pubsub := &mocks.PubSub{}
// 	event.On("Data").Return(data)
// 	data.On("BeforePublish", event).Return(nil)
// 	pubsub.On("Publish", event).Return(nil)
// 	world := &World{pubsub: pubsub}
// 	err := world.Publish(event)
// 	assert.NoError(suite.T(), err)
// 	data.AssertCalled(suite.T(), "BeforePublish", event)
// }
//
// func (suite *WorldTestSuite) TestPublishReturnsBeforePublishError() {
// 	err := errors.New("err")
// 	event := &mocks.Event{}
// 	data := &mocks.EventData{}
// 	pubsub := &mocks.PubSub{}
// 	event.On("Data").Return(data)
// 	data.On("BeforePublish", event).Return(err)
// 	world := &World{pubsub: pubsub}
// 	publishError := world.Publish(event)
// 	data.AssertCalled(suite.T(), "BeforePublish", event)
// 	assert.Equal(suite.T(), err, publishError)
// }
//
// func (suite *WorldTestSuite) TestPublishReturnsPubSubError() {
// 	err := errors.New("err")
// 	event := &mocks.Event{}
// 	data := &mocks.EventData{}
// 	pubsub := &mocks.PubSub{}
// 	event.On("Data").Return(data)
// 	data.On("BeforePublish", event).Return(nil)
// 	pubsub.On("Publish", event).Return(err)
// 	world := &World{pubsub: pubsub}
// 	publishError := world.Publish(event)
// 	assert.Equal(suite.T(), err, publishError)
// }
//
// func (suite *WorldTestSuite) TestNewWorldCallsSubscribe() {
// 	zone := &mocks.Zone{}
// 	suite.cache.On("Zone", ":0z").Return(zone, nil)
// 	suite.pubsub.On("Subscribe").Return(make(chan types.Event))
// 	_, err := newWorld("123", suite.cache, suite.factory, suite.pubsub, 2)
// 	assert.NoError(suite.T(), err)
// 	suite.pubsub.AssertCalled(suite.T(), "Subscribe")
// }
//
// // This is really gross :-(
// func (suite *WorldTestSuite) TestIncomingEventsCallOnReceive() {
// 	zone := &mocks.Zone{}
// 	suite.cache.On("Zone", ":0z").Return(zone, nil)
//
// 	ch := make(chan types.Event, 1)
//
// 	mockPubSub := &mocks.PubSub{}
// 	mockPubSub.On("Subscribe").Return(ch)
//
// 	_, err := newWorld("123", suite.cache, suite.factory, mockPubSub, 2)
// 	assert.NoError(suite.T(), err)
//
// 	mockEvent := &mocks.Event{}
// 	mockEventData := &mocks.EventData{}
//
// 	done := make(chan bool)
// 	mockEventData.On("OnReceive", mockEvent).Return(nil).Run(func(args mock.Arguments) {
// 		assert.Equal(suite.T(), mockEvent, args.Get(0))
// 		done <- true
// 	})
//
// 	mockEvent.On("Data").Return(mockEventData)
// 	mockEventData.AssertNotCalled(suite.T(), "OnReceive", mockEvent)
// 	ch <- mockEvent
// 	<-done
// 	mockPubSub.AssertCalled(suite.T(), "Subscribe")
// 	mockEventData.AssertCalled(suite.T(), "OnReceive", mockEvent)
// }
//
// func (suite *WorldTestSuite) TestSetZone() {
// 	zone := &mocks.Zone{}
// 	cache := &mocks.Cache{}
// 	cache.On("SetZone", zone).Return(nil)
// 	w := &World{cache: cache}
// 	err := w.SetZone(zone)
// 	assert.NoError(suite.T(), err)
// 	cache.AssertCalled(suite.T(), "SetZone", zone)
// }
//
// func (suite *WorldTestSuite) TestSetZoneReturnsError() {
// 	err := errors.New("sdfsd")
// 	zone := &mocks.Zone{}
// 	cache := &mocks.Cache{}
// 	cache.On("SetZone", zone).Return(err)
// 	w := &World{cache: cache}
// 	zerr := w.SetZone(zone)
// 	assert.Error(suite.T(), zerr)
// 	cache.AssertCalled(suite.T(), "SetZone", zone)
// 	assert.Equal(suite.T(), err, zerr)
// }
//
// func (suite *WorldTestSuite) TestSetUser() {
// 	user := &mocks.User{}
// 	cache := &mocks.Cache{}
// 	cache.On("SetUser", user).Return(nil)
// 	w := &World{cache: cache}
// 	err := w.SetUser(user)
// 	assert.NoError(suite.T(), err)
// 	cache.AssertCalled(suite.T(), "SetUser", user)
// }
//
// func (suite *WorldTestSuite) TestSetUserReturnsError() {
// 	err := errors.New("sdfsd")
// 	user := &mocks.User{}
// 	cache := &mocks.Cache{}
// 	cache.On("SetUser", user).Return(err)
// 	w := &World{cache: cache}
// 	zerr := w.SetUser(user)
// 	assert.Error(suite.T(), zerr)
// 	cache.AssertCalled(suite.T(), "SetUser", user)
// 	assert.Equal(suite.T(), err, zerr)
// }
//
// func TestWorldSuite(t *testing.T) {
// 	suite.Run(t, new(WorldTestSuite))
// }
