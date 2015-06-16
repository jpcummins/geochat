package chat

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EventTestSuite struct {
	suite.Suite
}

func (suite *EventTestSuite) SetupTest() {
}

func (suite *EventTestSuite) TestUnmarshalMessage() {
	event := &Event{}
	err := event.UnmarshalJSON(generateMockEventBytes("message"))
	assert.NoError(suite.T(), err)
}

func generateMockEventBytes(eventType string) []byte {
	return []byte("{\"type\":\"" + eventType + "\",\"id\":\"" + eventType + "_id\"}")
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
