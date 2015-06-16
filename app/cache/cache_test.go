package cache

import (
	"errors"
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CacheTestSuite struct {
	suite.Suite
	db    *mocks.DB
	cache *Cache
}

func (suite *CacheTestSuite) SetupTest() {
	suite.db = &mocks.DB{}
	suite.cache = NewCache(suite.db)
}

func (suite *CacheTestSuite) TestNewCache() {
	assert.NotNil(suite.T(), suite.cache.users)
	assert.NotNil(suite.T(), suite.cache.zones)
	assert.NotNil(suite.T(), suite.cache.db)
}

func (suite *CacheTestSuite) TestUserCallsDB() {
	user := &mocks.User{}
	suite.db.On("GetUser", "123").Return(user, nil)
	cachedUser, err := suite.cache.User("123")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user, cachedUser)
}

func (suite *CacheTestSuite) TestUserCallsDBAndReturnsError() {
	err := errors.New("bla")
	user := &mocks.User{}
	suite.db.On("GetUser", "123").Return(user, err)
	cachedUser, cachedErr := suite.cache.User("123")
	assert.Nil(suite.T(), cachedUser)
	assert.Equal(suite.T(), err, cachedErr)
}

func (suite *CacheTestSuite) TestUserRetrievesFromLocalCache() {
	mockUser := &mocks.User{}
	suite.cache.users["123"] = mockUser

	user, err := suite.cache.User("123")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, user)
	suite.db.AssertNotCalled(suite.T(), "GetUser", mock.Anything)
}

func (suite *CacheTestSuite) TestSetUserCachesAndCallsDB() {
	mockUser := &mocks.User{}
	mockUser.On("ID").Return("123")
	suite.db.On("SetUser", mockUser).Return(nil)
	err := suite.cache.SetUser(mockUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser, suite.cache.users["123"])
}

func (suite *CacheTestSuite) TestSetUserDBError() {
	err := errors.New("err")
	mockUser := &mocks.User{}
	mockUser.On("ID").Return("123")
	suite.db.On("SetUser", mockUser).Return(err)
	cachedError := suite.cache.SetUser(mockUser)
	assert.Error(suite.T(), cachedError)
	assert.Equal(suite.T(), err, cachedError)
	assert.Equal(suite.T(), 0, len(suite.cache.users))
}

func (suite *CacheTestSuite) TestZoneCallsDB() {
	zone := &mocks.Zone{}
	suite.db.On("GetZone", "123").Return(zone, nil)
	cachedZone, err := suite.cache.Zone("123")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), zone, cachedZone)
}

func (suite *CacheTestSuite) TestZoneCallsDBAndReturnsError() {
	err := errors.New("bla")
	zone := &mocks.Zone{}
	suite.db.On("GetZone", "123").Return(zone, err)
	cachedZone, cachedErr := suite.cache.Zone("123")
	assert.Nil(suite.T(), cachedZone)
	assert.Equal(suite.T(), err, cachedErr)
}

func (suite *CacheTestSuite) TestZoneRetrievesFromLocalCache() {
	mockZone := &mocks.Zone{}
	suite.cache.zones["123"] = mockZone

	zone, err := suite.cache.Zone("123")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockZone, zone)
	suite.db.AssertNotCalled(suite.T(), "GetZone", mock.Anything)
}

func (suite *CacheTestSuite) TestSetZoneCachesAndCallsDB() {
	mockZone := &mocks.Zone{}
	mockZone.On("ID").Return("123")
	suite.db.On("SetZone", mockZone).Return(nil)
	err := suite.cache.SetZone(mockZone)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockZone, suite.cache.zones["123"])
}

func (suite *CacheTestSuite) TestSetZoneDBError() {
	err := errors.New("err")
	mockZone := &mocks.Zone{}
	mockZone.On("ID").Return("123")
	suite.db.On("SetZone", mockZone).Return(err)
	cachedError := suite.cache.SetZone(mockZone)
	assert.Error(suite.T(), cachedError)
	assert.Equal(suite.T(), err, cachedError)
	assert.Equal(suite.T(), 0, len(suite.cache.zones))
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
