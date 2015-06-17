package chat

import (
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EventTestSuite struct {
	suite.Suite
	mockChat *mocks.Chat
}

func (suite *EventTestSuite) SetupTest() {
	suite.mockChat = &mocks.Chat{}
	chat = suite.mockChat

	world := &mocks.World{}
	suite.mockChat.On("World", "wid").Return(world, nil)
}

func (suite *EventTestSuite) TestNewEvent() {
	data := &mocks.EventData{}
	data.On("Type").Return("test")

	e, err := newEvent("eventid", "wid", data)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "eventid", e.ID())
	assert.Equal(suite.T(), "test", e.Type())
	assert.Equal(suite.T(), data, e.Data())
	assert.Equal(suite.T(), "wid", e.eventJSON.WorldID)
}

func (suite *EventTestSuite) TestUnmarshalMessage() {
	event := &Event{}
	err := event.UnmarshalJSON(generateMockEventBytes("message"))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "eventid", event.ID())
	assert.Equal(suite.T(), "message", event.Type())
	assert.Equal(suite.T(), "wid", event.eventJSON.WorldID)
}

func (suite *EventTestSuite) TestEventUnmarshalError() {

	event := &Event{}
	err := event.UnmarshalJSON(generateMockEventBytes("bad"))
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "Unable to unmarshal command: bad", err.Error())
}

func generateMockEventBytes(eventType string) []byte {
	return []byte("{\"id\":\"eventid\",\"type\":\"" + eventType + "\",\"world_id\":\"wid\",\"data\":{}}")
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
