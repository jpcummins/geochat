package chat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var lat float64 = 47.6235616
var long float64 = -122.330341

func TestFindChatZone_EmptyWorld(t *testing.T) {
	root := createZone("", '0', 'z', nil)
	zone, _ := findChatZone(root, "c23nb")
	assert.Equal(t, root, zone)
	assert.Equal(t, ":0g", zone.left.Zonehash)
	assert.Equal(t, ":hz", zone.right.Zonehash)
}

func TestFindChatZone_SubhashToFirstGeohash3(t *testing.T) {
	testCases := []string{"000", "z0z", "2k1", "bbc", "zzz", "c23nb"}

	var root *Zone
	var zone *Zone

	for _, test := range testCases {
		root = createZone("", '0', 'z', nil)
		root.Count = maxRoomSize
		for i := 0; i < 5*len(test); i++ {
			zone, _ = findChatZone(root, test)
			zone.Count = maxRoomSize
			if len(zone.geohash) == len(test) {
				break
			}
		}
		assert.Equal(t, test+":0z", zone.Zonehash)
	}
}
