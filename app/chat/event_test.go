package chat

import (
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EventTestSuite struct {
	suite.Suite
}

func (suite *EventTestSuite) TestNewEvent() {
	world := &mocks.World{}
	world.On("ID").Return("worldid")

	data := &mocks.EventData{}
	data.On("Type").Return("test")

	e := newEvent("eventid", world, data)
	assert.Equal(suite.T(), "eventid", e.ID())
	assert.Equal(suite.T(), "test", e.Type())
	assert.Equal(suite.T(), data, e.Data())
}

func (suite *EventTestSuite) TestUnmarshalMessage() {
	event := &Event{}
	err := event.UnmarshalJSON(generateMockEventBytes("message"))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "eventid", event.ID())
	assert.Equal(suite.T(), "message", event.Type())
	assert.Equal(suite.T(), "worldid", event.WorldID())
}

func (suite *EventTestSuite) TestEventUnmarshalError() {
	event := &Event{}
	err := event.UnmarshalJSON(generateMockEventBytes("bad"))
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "Unable to unmarshal command: bad", err.Error())
}

func generateMockEventBytes(eventType string) []byte {
	return []byte("{\"id\":\"eventid\",\"type\":\"" + eventType + "\",\"world_id\":\"worldid\",\"data\":{}}")
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
