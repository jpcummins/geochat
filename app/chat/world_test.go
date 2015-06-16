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

	world, err := newWorld("", suite.cache, suite.factory, 2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockZone, world.root)
	assert.Equal(suite.T(), suite.cache, world.cache)
	assert.Equal(suite.T(), 2, world.maxUsersPerZone)
}

func (suite *WorldTestSuite) TestNewWorldReturnsError() {
	worldErr := errors.New("err")
	suite.cache.On("Zone", ":0z").Return(nil, worldErr)
	world, err := newWorld("", suite.cache, suite.factory, 2)
	assert.Equal(suite.T(), err, worldErr)
	assert.Nil(suite.T(), world)
}

func (suite *WorldTestSuite) TestNewWorldErrorOnCreation() {
	err := errors.New("test error")
	suite.cache.On("Zone", ":0z").Return(nil, err)

	world, err := newWorld("", suite.cache, suite.factory, 2)
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
		id:              "test",
		cache:           suite.cache,
		factory:         suite.factory,
		maxUsersPerZone: 22,
	}

	suite.factory.On("NewZone", ":0z", "test", 22).Return(zone, nil)

	z, err := world.GetOrCreateZone(":0z")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), zone, z)

	suite.cache.AssertCalled(suite.T(), "Zone", ":0z")
	suite.cache.AssertCalled(suite.T(), "SetZone", zone)
	suite.factory.AssertCalled(suite.T(), "NewZone", ":0z", "test", 22)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMissAndNewZoneError() {
	newZoneErr := errors.New("err")
	suite.cache.On("Zone", ":0z").Return(nil, nil)

	world := &World{
		id:              "test",
		cache:           suite.cache,
		factory:         suite.factory,
		maxUsersPerZone: 22,
	}

	suite.factory.On("NewZone", ":0z", "test", 22).Return(nil, newZoneErr)

	z, err := world.GetOrCreateZone(":0z")
	assert.Nil(suite.T(), z)
	assert.Equal(suite.T(), newZoneErr, err)

	suite.cache.AssertCalled(suite.T(), "Zone", ":0z")
	suite.factory.AssertCalled(suite.T(), "NewZone", ":0z", "test", 22)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneCacheMissAndSetZoneError() {
	zone := &mocks.Zone{}
	setZoneErr := errors.New("err")
	suite.cache.On("Zone", ":0z").Return(nil, nil)
	suite.cache.On("SetZone", zone).Return(setZoneErr)

	world := &World{
		id:              "test",
		cache:           suite.cache,
		factory:         suite.factory,
		maxUsersPerZone: 22,
	}

	suite.factory.On("NewZone", ":0z", "test", 22).Return(zone, nil)

	z, err := world.GetOrCreateZone(":0z")
	assert.Nil(suite.T(), z)
	assert.Equal(suite.T(), setZoneErr, err)

	suite.cache.AssertCalled(suite.T(), "Zone", ":0z")
	suite.cache.AssertCalled(suite.T(), "SetZone", zone)
	suite.factory.AssertCalled(suite.T(), "NewZone", ":0z", "test", 22)
}

func (suite *WorldTestSuite) TestMultipleWorldsWithSameDBDependencyReturnsSameRoot() {
	zone := &mocks.Zone{}
	suite.cache.On("Zone", ":0z").Return(zone, nil)

	world1, err1 := newWorld("", suite.cache, nil, 1)
	world2, err2 := newWorld("", suite.cache, nil, 1)
	world3, err3 := newWorld("", suite.cache, nil, 1)

	assert.Equal(suite.T(), zone, world1.root)
	assert.Equal(suite.T(), zone, world2.root)
	assert.Equal(suite.T(), zone, world3.root)

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NoError(suite.T(), err3)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_EmptyWorld() {
	root := &mocks.Zone{}
	root.On("IsOpen").Return(true)

	user := &mocks.User{}
	user.On("Location").Return(seattle)

	world := &World{
		cache:           suite.cache,
		root:            root,
		maxUsersPerZone: 2,
	}

	zone, err := world.GetOrCreateZoneForUser(user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), root, zone)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_RightZoneReturnsError() {
	err := errors.New("err")

	root := &mocks.Zone{}
	root.On("IsOpen").Return(false)
	root.On("Geohash").Return("")
	root.On("RightZoneID").Return(":hz")

	suite.cache.On("Zone", ":hz").Return(nil, err)

	user := &mocks.User{}
	user.On("Location").Return(seattle)

	world := &World{
		cache:           suite.cache,
		root:            root,
		maxUsersPerZone: 2,
	}

	zone, err := world.GetOrCreateZoneForUser(user)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), zone)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_LeftZoneReturnsError() {
	err := errors.New("err")

	root := &mocks.Zone{}
	root.On("IsOpen").Return(false)
	root.On("Geohash").Return("")
	root.On("RightZoneID").Return(":hz")
	root.On("LeftZoneID").Return(":0g")

	suite.cache.On("Zone", ":hz").Return(&mocks.Zone{}, nil)
	suite.cache.On("Zone", ":0g").Return(nil, err)

	user := &mocks.User{}
	user.On("Location").Return(seattle)

	world := &World{
		cache:           suite.cache,
		root:            root,
		maxUsersPerZone: 2,
	}

	zone, err := world.GetOrCreateZoneForUser(user)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), zone)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_ReturnsLeftZone() {
	root := &mocks.Zone{}
	root.On("IsOpen").Return(false)
	root.On("Geohash").Return("")
	root.On("RightZoneID").Return(":hz")
	root.On("LeftZoneID").Return(":0g")

	leftZone := &mocks.Zone{}
	leftZone.On("IsOpen").Return(true)

	rightZone := &mocks.Zone{}
	rightZone.On("Geohash").Return("")
	rightZone.On("From").Return("h")

	suite.cache.On("Zone", ":0g").Return(leftZone, nil)
	suite.cache.On("Zone", ":hz").Return(rightZone, nil)

	user := &mocks.User{}
	user.On("Location").Return(seattle)

	world := &World{
		cache:           suite.cache,
		root:            root,
		maxUsersPerZone: 2,
	}

	zone, err := world.GetOrCreateZoneForUser(user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), leftZone, zone)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_ReturnsRightZone() {
	root := &mocks.Zone{}
	root.On("IsOpen").Return(false)
	root.On("Geohash").Return("")
	root.On("RightZoneID").Return(":hz")
	root.On("LeftZoneID").Return(":0g")

	leftZone := &mocks.Zone{}

	rightZone := &mocks.Zone{}
	rightZone.On("Geohash").Return("")
	rightZone.On("From").Return("h")
	rightZone.On("IsOpen").Return(true)

	suite.cache.On("Zone", ":0g").Return(leftZone, nil)
	suite.cache.On("Zone", ":hz").Return(rightZone, nil)

	user := &mocks.User{}
	user.On("Location").Return(rome)

	world := &World{
		cache:           suite.cache,
		root:            root,
		maxUsersPerZone: 2,
	}

	zone, err := world.GetOrCreateZoneForUser(user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), rightZone, zone)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_ReturnsLeftZone2() {
	root := &mocks.Zone{}
	root.On("IsOpen").Return(false)
	root.On("Geohash").Return("")
	root.On("RightZoneID").Return(":hz")
	root.On("LeftZoneID").Return(":0g")

	leftZone := &mocks.Zone{}
	leftZone.On("IsOpen").Return(true)

	rightZone := &mocks.Zone{}
	rightZone.On("Geohash").Return("s")
	rightZone.On("From").Return("h")

	suite.cache.On("Zone", ":0g").Return(leftZone, nil)
	suite.cache.On("Zone", ":hz").Return(rightZone, nil)

	user := &mocks.User{}
	user.On("Location").Return(seattle)

	world := &World{
		cache:           suite.cache,
		root:            root,
		maxUsersPerZone: 2,
	}

	zone, err := world.GetOrCreateZoneForUser(user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), leftZone, zone)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_ReturnsRightZone2() {
	root := &mocks.Zone{}
	root.On("IsOpen").Return(false)
	root.On("Geohash").Return("")
	root.On("RightZoneID").Return(":hz")
	root.On("LeftZoneID").Return(":0g")

	leftZone := &mocks.Zone{}

	rightZone := &mocks.Zone{}
	rightZone.On("Geohash").Return("s")
	rightZone.On("From").Return("h")
	rightZone.On("IsOpen").Return(true)

	suite.cache.On("Zone", ":0g").Return(leftZone, nil)
	suite.cache.On("Zone", ":hz").Return(rightZone, nil)

	user := &mocks.User{}
	user.On("Location").Return(rome)

	world := &World{
		cache:           suite.cache,
		root:            root,
		maxUsersPerZone: 2,
	}

	zone, err := world.GetOrCreateZoneForUser(user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), rightZone, zone)
}

func (suite *WorldTestSuite) TestGetOrCreateZoneForUser_ErrorOnNoOpenRooms() {
	root := &mocks.Zone{}
	root.On("IsOpen").Return(false)
	root.On("Geohash").Return(seattle.Geohash())

	user := &mocks.User{}
	user.On("Location").Return(seattle)

	world := &World{
		cache:           suite.cache,
		root:            root,
		maxUsersPerZone: 2,
	}

	zone, err := world.GetOrCreateZoneForUser(user)
	assert.Nil(suite.T(), zone)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "Unable to find zone", err.Error())
}

func (suite *WorldTestSuite) TestIntegration() {
	db := &mocks.DB{}
	db.On("GetZone", mock.Anything).Return(nil, nil)
	db.On("SetZone", mock.Anything).Return(nil)

	world := &World{
		id:              "test",
		cache:           cache.NewCache(db),
		factory:         &Factory{},
		maxUsersPerZone: 1,
	}

	w, err := world.GetOrCreateZone(":0z")
	assert.NoError(suite.T(), err)
	world.root = w
	world.root.SetIsOpen(false)
	assert.Equal(suite.T(), ":0z", world.root.ID())

	user := &mocks.User{}
	user.On("Location").Return(seattle)

	zone, err := world.GetOrCreateZoneForUser(user)
	assert.Equal(suite.T(), ":0g", zone.ID())

	zone.SetIsOpen(false)
	zone, err = world.GetOrCreateZoneForUser(user)
	assert.Equal(suite.T(), ":8g", zone.ID())
}

func (suite *WorldTestSuite) TestIntegration2() {
	testCases := []string{"000", "z0z", "2k1", "bbc", "zzz", "c23nb"}

	db := &mocks.DB{}
	db.On("GetZone", mock.Anything).Return(nil, nil)
	db.On("SetZone", mock.Anything).Return(nil)

	for _, test := range testCases {
		world := &World{
			cache:           cache.NewCache(db),
			factory:         &Factory{},
			maxUsersPerZone: 1,
		}
		w, err := world.GetOrCreateZone(":0z")
		assert.NoError(suite.T(), err)
		world.root = w
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
		assert.Equal(suite.T(), test+":0z", zone.ID())
	}
}

func TestWorldSuite(t *testing.T) {
	suite.Run(t, new(WorldTestSuite))
}
