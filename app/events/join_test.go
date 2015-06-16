package events

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JoinTestSuite struct {
	suite.Suite
}

func (suite *JoinTestSuite) SetupTest() {

}

func (suite *JoinTestSuite) TestNewJoinEvent() {
	j, err := JoinEventData("123")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "123", j.userID)
}

func TestJointSuite(t *testing.T) {
	suite.Run(t, new(JoinTestSuite))
}
