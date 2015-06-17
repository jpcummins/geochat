package chat

import (
	"errors"
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EventTestSuite struct {
	suite.Suite
	chat  *mocks.Chat
	world *mocks.World
}

func (suite *EventTestSuite) SetupTest() {
	suite.chat = &mocks.Chat{}
	chat = suite.chat
	suite.world = &mocks.World{}
	suite.chat.On("World", "wid").Return(suite.world, nil)
	suite.world.On("ID").Return("wid")
}

func (suite *EventTestSuite) TestNewEvent() {
	data := &mocks.EventData{}
	data.On("Type").Return("test")
	e, err := newEvent("eventid", suite.world.ID(), data)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "eventid", e.ID())
	assert.Equal(suite.T(), "test", e.Type())
	assert.Equal(suite.T(), data, e.Data())
	assert.Equal(suite.T(), suite.world, e.World())
	assert.Equal(suite.T(), suite.world.ID(), e.eventJSON.WorldID)
}

func (suite *EventTestSuite) TestNewEventError() {
	err := errors.New("error")
	suite.chat.On("World", "error_world").Return(nil, err)
	testE, testErr := newEvent("eventid", "error_world", &mocks.EventData{})
	assert.Equal(suite.T(), err, testErr)
	assert.Nil(suite.T(), testE)
}

func (suite *EventTestSuite) TestNewEventError2() {
	suite.chat.On("World", "error_world").Return(nil, nil)
	testE, testErr := newEvent("eventid", "error_world", &mocks.EventData{})
	assert.Equal(suite.T(), "Unable to find world: error_world", testErr.Error())
	assert.Nil(suite.T(), testE)
}

func (suite *EventTestSuite) TestUnmarshalMessage() {
	event := &Event{}
	err := event.UnmarshalJSON(generateMockEventBytes("message", "wid"))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "eventid", event.ID())
	assert.Equal(suite.T(), "message", event.Type())
	assert.Equal(suite.T(), suite.world, event.World())
	assert.Equal(suite.T(), "wid", event.eventJSON.WorldID)
}

func (suite *EventTestSuite) TestEventUnmarshalError() {
	event := &Event{}
	err := event.UnmarshalJSON(generateMockEventBytes("bad", "wid"))
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "Unable to unmarshal command: bad", err.Error())
}

func (suite *EventTestSuite) TestEventUnmarshalError2() {
	err := errors.New("error")
	suite.chat.On("World", "error_world").Return(nil, err)
	event := &Event{}
	testErr := event.UnmarshalJSON(generateMockEventBytes("bad", "error_world"))
	assert.Equal(suite.T(), err, testErr)
}

func (suite *EventTestSuite) TestEventUnmarshalError3() {
	suite.chat.On("World", "error_world").Return(nil, nil)
	event := &Event{}
	testErr := event.UnmarshalJSON(generateMockEventBytes("bad", "error_world"))
	assert.Equal(suite.T(), "Unable to find world: error_world", testErr.Error())
}

func generateMockEventBytes(eventType string, worldID string) []byte {
	return []byte("{\"id\":\"eventid\",\"type\":\"" + eventType + "\",\"world_id\":\"" + worldID + "\",\"data\":{}}")
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
