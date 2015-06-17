package events

import (
	"errors"
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JoinTestSuite struct {
	suite.Suite
}

func (suite *JoinTestSuite) TestNewJoinEvent() {
	world := &mocks.World{}
	zone := &mocks.Zone{}
	user := &mocks.User{}

	j, err := NewJoin(world, zone, user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), world, j.world)
	assert.Equal(suite.T(), zone, j.zone)
	assert.Equal(suite.T(), user, j.user)
}

func (suite *JoinTestSuite) TestBeforePublishSavesUser() {
	world := &mocks.World{}
	zone := &mocks.Zone{}
	user := &mocks.User{}

	world.On("SetUser", user).Return(nil)
	world.On("SetZone", zone).Return(nil)
	zone.On("AddUser", user).Return()

	j, err := NewJoin(world, zone, user)
	j.BeforePublish(&mocks.Event{})
	assert.NoError(suite.T(), err)
	world.AssertCalled(suite.T(), "SetUser", user)
}

func (suite *JoinTestSuite) TestBeforePublishErrors() {
	mockError := errors.New("err")

	world := &mocks.World{}
	zone := &mocks.Zone{}
	user := &mocks.User{}

	world.On("SetUser", user).Return(mockError)
	world.On("SetZone", zone).Return(nil)

	j, err := NewJoin(world, zone, user)
	err = j.BeforePublish(&mocks.Event{})
	assert.Error(suite.T(), err)
	world.AssertCalled(suite.T(), "SetUser", user)
	assert.Equal(suite.T(), mockError, err)
}

func (suite *JoinTestSuite) TestBeforePublishErrors2() {
	mockError := errors.New("err")

	world := &mocks.World{}
	zone := &mocks.Zone{}
	user := &mocks.User{}

	world.On("SetUser", user).Return(nil)
	zone.On("AddUser", user).Return()
	world.On("SetZone", zone).Return(mockError)

	j, err := NewJoin(world, zone, user)
	err = j.BeforePublish(&mocks.Event{})
	assert.Error(suite.T(), err)
	world.AssertCalled(suite.T(), "SetZone", zone)
	assert.Equal(suite.T(), mockError, err)
}

func TestJointSuite(t *testing.T) {
	suite.Run(t, new(JoinTestSuite))
}
