package pubsub

import (
	"github.com/jpcummins/geochat/app/mocks"
	"github.com/jpcummins/geochat/app/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EventTestSuite struct {
	suite.Suite
	world *mocks.World
}

func (suite *EventTestSuite) SetupTest() {
	suite.world = &mocks.World{}
	suite.world.On("ID").Return("wid")
}

func (suite *EventTestSuite) TestNewEvent() {
	data := &mocks.ServerEventData{}
	eventType := types.ServerEventType("test")
	data.On("Type").Return(eventType)
	e := newServerEvent("eventid", suite.world, data)
	assert.Equal(suite.T(), "eventid", e.ID())
	assert.Equal(suite.T(), eventType, e.Type())
	assert.Equal(suite.T(), data, e.Data())
	assert.Equal(suite.T(), suite.world, e.World())
}

func (suite *EventTestSuite) TestUnmarshalMessage() {
	event := &ServerEvent{}
	eventType := types.ServerEventType("message")
	err := event.UnmarshalJSON(generateMockEventBytes(eventType, "wid"))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "eventid", event.ID())
	assert.Equal(suite.T(), eventType, event.Type())
}

func (suite *EventTestSuite) TestEventUnmarshalError() {
	event := &ServerEvent{}
	err := event.UnmarshalJSON(generateMockEventBytes("bad", "wid"))
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "Unable to unmarshal command: bad", err.Error())
}

func generateMockEventBytes(eventType types.ServerEventType, worldID string) []byte {
	return []byte("{\"id\":\"eventid\",\"type\":\"" + string(eventType) + "\",\"world_id\":\"" + worldID + "\",\"data\":{}}")
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
