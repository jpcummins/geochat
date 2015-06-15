package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type WorldTestSuite struct {
	suite.Suite
	cache   *mocks.Cache
	factory *mocks.Factory
}

func (suite *WorldTestSuite) SetupTest() {
	suite.cache = &mocks.Cache{}
	suite.factory = &mocks.Factory{}
}

func (suite *WorldTestSuite) TestNewWorld() {
	mockZone := &mocks.Zone{}
	suite.cache.On("Zone", ":0z").Return(mockZone, nil)

	world, err := newWorld(suite.cache, suite.factory, 2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockZone, world.root)
	assert.Equal(suite.T(), suite.cache, world.cache)
	assert.Equal(suite.T(), 2, world.maxUsersPerZone)
}

func (suite *WorldTestSuite) TestNewWorldReturnsError() {
	worldErr := errors.New("err")
	suite.cache.On("Zone", ":0z").Return(nil, worldErr)
	world, err := newWorld(suite.cache, suite.factory, 2)
	assert.Equal(suite.T(), err, worldErr)
	assert.Nil(suite.T(), world)
}

func (suite *WorldTestSuite) TestNewWorldErrorOnCreation() {
	err := errors.New("test error")
	suite.cache.On("Zone", ":0z").Return(nil, err)

	world, err := newWorld(suite.cache, suite.factory, 2)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), world)
}

func (suite *WorldTestSuite) TestGetOrCreateZone() {
	world := &World{
		cache:           suite.cache,
		maxUsersPerZone: 1,
	}

	zone := &mocks.Zone{}
	suite.cache.On("Zone", ":0z").Return(zone, nil)

	z, err := world.GetOrCreateZone(":0z")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), zone, z)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMiss() {
	zone := &mocks.Zone{}
	suite.cache.On("Zone", ":0z").Return(nil, nil)
	suite.cache.On("SetZone", zone).Return(nil)

	world := &World{
		cache:           suite.cache,
		factory:         suite.factory,
		maxUsersPerZone: 22,
	}

	suite.factory.On("NewZone", world, ":0z").Return(zone, nil)

	z, err := world.GetOrCreateZone(":0z")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), zone, z)

	suite.cache.AssertCalled(suite.T(), "Zone", ":0z")
	suite.cache.AssertCalled(suite.T(), "SetZone", zone)
	suite.factory.AssertCalled(suite.T(), "NewZone", world, ":0z")
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMissAndNewZoneError() {
	newZoneErr := errors.New("err")
	suite.cache.On("Zone", ":0z").Return(nil, nil)

	world := &World{
		cache:           suite.cache,
		factory:         suite.factory,
		maxUsersPerZone: 22,
	}

	suite.factory.On("NewZone", world, ":0z").Return(nil, newZoneErr)

	z, err := world.GetOrCreateZone(":0z")
	assert.Nil(suite.T(), z)
	assert.Equal(suite.T(), newZoneErr, err)

	suite.cache.AssertCalled(suite.T(), "Zone", ":0z")
	suite.factory.AssertCalled(suite.T(), "NewZone", world, ":0z")
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMissAndSetZoneError() {
	zone := &mocks.Zone{}
	setZoneErr := errors.New("err")
	suite.cache.On("Zone", ":0z").Return(nil, nil)
	suite.cache.On("SetZone", zone).Return(setZoneErr)

	world := &World{
		cache:           suite.cache,
		factory:         suite.factory,
		maxUsersPerZone: 22,
	}

	suite.factory.On("NewZone", world, ":0z").Return(zone, nil)

	z, err := world.GetOrCreateZone(":0z")
	assert.Nil(suite.T(), z)
	assert.Equal(suite.T(), setZoneErr, err)

	suite.cache.AssertCalled(suite.T(), "Zone", ":0z")
	suite.cache.AssertCalled(suite.T(), "SetZone", zone)
	suite.factory.AssertCalled(suite.T(), "NewZone", world, ":0z")
}

func (suite *WorldTestSuite) TestMultipleWorldsWithSameDBDependencyReturnsSameRoot() {
	zone := &mocks.Zone{}
	suite.cache.On("Zone", ":0z").Return(zone, nil)

	world1, err1 := newWorld(suite.cache, nil, 1)
	world2, err2 := newWorld(suite.cache, nil, 1)
	world3, err3 := newWorld(suite.cache, nil, 1)

	assert.Equal(suite.T(), zone, world1.root)
	assert.Equal(suite.T(), zone, world2.root)
	assert.Equal(suite.T(), zone, world3.root)

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NoError(suite.T(), err3)
}

func TestWorldSuite(t *testing.T) {
	suite.Run(t, new(WorldTestSuite))
}

//
// import (
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"testing"
// )
//
// type ZoneTestSuite struct {
// 	suite.Suite
// 	Zone *Zone
// }
//
// func (suite *ZoneTestSuite) SetupTest() {
// 	suite.Zone = newZone("", '0', 'z', nil, 1)
// 	connection = &mockConnection{}
// }
//
// var lat float64 = 47.6235616
// var long float64 = -122.330341
//
// func (suite *ZoneTestSuite) TestFindChatZone_EmptyWorld() {
// 	root := suite.Zone
// 	zone, _ := findChatZone(root, "c23nb")
// 	assert.Equal(suite.T(), root, zone)
// 	assert.Equal(suite.T(), ":0g", zone.left.id)
// 	assert.Equal(suite.T(), ":hz", zone.right.id)
// }
//
// func (suite *ZoneTestSuite) TestFindChatZone_ZoneCreation() {
// 	testCases := []string{"000", "z0z", "2k1", "bbc", "zzz", "c23nb"}
//
// 	var root *Zone
// 	var zone *Zone
//
// 	for _, test := range testCases {
// 		root = newZone("", '0', 'z', nil, 1)
// 		root.isOpen = false
// 		for i := 0; i < 5*len(test); i++ {
// 			zone, _ = findChatZone(root, test)
// 			zone.isOpen = false
// 			if len(zone.geohash) == len(test) {
// 				break
// 			}
// 		}
// 		assert.Equal(suite.T(), test+":0z", zone.id)
//
// 		// While we're at it, test GetZone() functionality.
// 		world = &World{root, nil, nil}
// 		z, err := getOrCreateZone(zone.id)
// 		assert.NoError(suite.T(), err)
// 		assert.Equal(suite.T(), zone, z)
// 	}
// }
//
// func (suite *ZoneTestSuite) TestGetZone() {
// 	root := newZone("", '0', 'z', nil, 1)
// 	world = &World{root, nil, nil}
// 	zone, err := getOrCreateZone(":0z")
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), world.root, zone)
// }
//
// func (suite *ZoneTestSuite) TestZoneSplit() {
// 	root := newZone("", '0', 'z', nil, 4)
//
// 	users := make([]*User, 5)
// 	for i, _ := range users {
// 		users[i] = NewUser(lat, long, "user")
// 		root.join(users[i])
//
// 		assert.True(suite.T(), root.isOpen)
// 	}
//
// }
//
// func TestSuite(t *testing.T) {
// 	suite.Run(t, new(ZoneTestSuite))
// }
