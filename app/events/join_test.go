package events

import (
	"errors"
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JoinTestSuite struct {
	suite.Suite
	world *mocks.World
	zone  *mocks.Zone
	user  *mocks.User
	users *mocks.Users
	zones *mocks.Zones
	event *mocks.Event
}

func (suite *JoinTestSuite) SetupTest() {
	suite.world = &mocks.World{}
	suite.zone = &mocks.Zone{}
	suite.user = &mocks.User{}
	suite.users = &mocks.Users{}
	suite.zones = &mocks.Zones{}
	suite.event = &mocks.Event{}

	suite.zone.On("ID").Return("zoneid")
	suite.user.On("ID").Return("userid")
	suite.user.On("Name").Return("username")
	suite.event.On("World").Return(suite.world)
	suite.world.On("Zone").Return(suite.zone)
	suite.world.On("Users").Return(suite.users)
	suite.world.On("Zones").Return(suite.zones)
}

func (suite *JoinTestSuite) TestNewJoinEvent() {
	j, err := NewJoin(suite.zone, suite.user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), suite.zone, j.zone)
	assert.Equal(suite.T(), suite.user, j.user)
}

func (suite *JoinTestSuite) TestBeforePublishSavesUser() {
	suite.users.On("SetUser", suite.user).Return(nil)
	suite.zones.On("SetZone", suite.zone).Return(nil)
	suite.zone.On("AddUser", suite.user).Return()

	j, _ := NewJoin(suite.zone, suite.user)
	err := j.BeforePublish(suite.event)
	assert.NoError(suite.T(), err)
	suite.users.AssertCalled(suite.T(), "SetUser", suite.user)
}

func (suite *JoinTestSuite) TestBeforePublishErrors() {
	err1 := errors.New("dflksdj")
	suite.users.On("SetUser", suite.user).Return(err1)

	j, _ := NewJoin(suite.zone, suite.user)
	err2 := j.BeforePublish(suite.event)
	assert.Equal(suite.T(), err1, err2)
}

func (suite *JoinTestSuite) TestBeforePublishErrors2() {
	err1 := errors.New("sdflksdjf")
	suite.users.On("SetUser", suite.user).Return(nil)
	suite.zones.On("SetZone", suite.zone).Return(err1)
	suite.zone.On("AddUser", suite.user).Return()

	j, _ := NewJoin(suite.zone, suite.user)
	err2 := j.BeforePublish(suite.event)
	assert.Equal(suite.T(), err1, err2)
}

func TestJointSuite(t *testing.T) {
	suite.Run(t, new(JoinTestSuite))
}
