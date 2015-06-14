package chat

import (
	"github.com/jpcummins/geochat/app/cache"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type WorldTestSuite struct {
	suite.Suite
	cache *cache.MockCache
	world *World
}

func (suite *WorldTestSuite) SetupTest() {
	suite.cache = &cache.MockCache{}
}

func (suite *WorldTestSuite) TestNewZone() {

}

func TestWorldSuite(t *testing.T) {
	suite.Run(t, new(WorldTestSuite))
}

//
// import (
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"testing"
// )
//
// type ZoneTestSuite struct {
// 	suite.Suite
// 	Zone *Zone
// }
//
// func (suite *ZoneTestSuite) SetupTest() {
// 	suite.Zone = newZone("", '0', 'z', nil, 1)
// 	connection = &mockConnection{}
// }
//
// var lat float64 = 47.6235616
// var long float64 = -122.330341
//
// func (suite *ZoneTestSuite) TestFindChatZone_EmptyWorld() {
// 	root := suite.Zone
// 	zone, _ := findChatZone(root, "c23nb")
// 	assert.Equal(suite.T(), root, zone)
// 	assert.Equal(suite.T(), ":0g", zone.left.id)
// 	assert.Equal(suite.T(), ":hz", zone.right.id)
// }
//
// func (suite *ZoneTestSuite) TestFindChatZone_ZoneCreation() {
// 	testCases := []string{"000", "z0z", "2k1", "bbc", "zzz", "c23nb"}
//
// 	var root *Zone
// 	var zone *Zone
//
// 	for _, test := range testCases {
// 		root = newZone("", '0', 'z', nil, 1)
// 		root.isOpen = false
// 		for i := 0; i < 5*len(test); i++ {
// 			zone, _ = findChatZone(root, test)
// 			zone.isOpen = false
// 			if len(zone.geohash) == len(test) {
// 				break
// 			}
// 		}
// 		assert.Equal(suite.T(), test+":0z", zone.id)
//
// 		// While we're at it, test GetZone() functionality.
// 		world = &World{root, nil, nil}
// 		z, err := getOrCreateZone(zone.id)
// 		assert.NoError(suite.T(), err)
// 		assert.Equal(suite.T(), zone, z)
// 	}
// }
//
// func (suite *ZoneTestSuite) TestGetZone() {
// 	root := newZone("", '0', 'z', nil, 1)
// 	world = &World{root, nil, nil}
// 	zone, err := getOrCreateZone(":0z")
// 	assert.NoError(suite.T(), err)
// 	assert.Equal(suite.T(), world.root, zone)
// }
//
// func (suite *ZoneTestSuite) TestZoneSplit() {
// 	root := newZone("", '0', 'z', nil, 4)
//
// 	users := make([]*User, 5)
// 	for i, _ := range users {
// 		users[i] = NewUser(lat, long, "user")
// 		root.join(users[i])
//
// 		assert.True(suite.T(), root.isOpen)
// 	}
//
// }
//
// func TestSuite(t *testing.T) {
// 	suite.Run(t, new(ZoneTestSuite))
// }
