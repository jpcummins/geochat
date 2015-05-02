package chat

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ZoneTestSuite struct {
	suite.Suite
	Zone *Zone
}

func (suite *ZoneTestSuite) SetupTest() {
	suite.Zone = newZone("", '0', 'z', nil, 1)
}

var lat float64 = 47.6235616
var long float64 = -122.330341

func (suite *ZoneTestSuite) TestFindChatZone_EmptyWorld() {
	root := suite.Zone
	zone, _ := findChatZone(root, "c23nb")
	assert.Equal(suite.T(), root, zone)
	assert.Equal(suite.T(), ":0g", zone.left.Zonehash)
	assert.Equal(suite.T(), ":hz", zone.right.Zonehash)
}

func (suite *ZoneTestSuite) TestFindChatZone_ZoneCreation() {
	testCases := []string{"000", "z0z", "2k1", "bbc", "zzz", "c23nb"}

	var root *Zone
	var zone *Zone

	for _, test := range testCases {
		root = newZone("", '0', 'z', nil, 1)
		root.count = 1
		for i := 0; i < 5*len(test); i++ {
			zone, _ = findChatZone(root, test)
			zone.count = 1
			if len(zone.geohash) == len(test) {
				break
			}
		}
		assert.Equal(suite.T(), test+":0z", zone.Zonehash)

		// While we're at it, test GetZone() functionality.
		world = &World{root, nil, nil}
		z, err := getOrCreateZone(zone.Zonehash)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), zone, z)
	}
}

func (suite *ZoneTestSuite) TestGetZone() {
	world = &World{suite.Zone, nil, nil}
	zone, err := getOrCreateZone(":0z")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), world.root, zone)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ZoneTestSuite))
}
