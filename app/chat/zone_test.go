package chat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var lat  float64 = 47.6235616
var long float64 = -122.330341

func TestFindChatZone_EmptyWorld(t *testing.T) {
	zone, _ := FindAvailableZone(lat, long)
	assert.Equal(t, world, zone)
	assert.Equal(t, ":2", world.left.Zonehash)
	assert.Equal(t, ":3", world.right.Zonehash)
}

func TestFindChatZone_SubhashToFirstGeohash(t *testing.T) {
	root := createZone("", 1, nil)
	root.Count = maxRoomSize
	assert.Equal(t, ":1", root.Zonehash)

	zone, _ := findChatZone(root, "c23nb")
	zone.Count = maxRoomSize
	assert.Equal(t, ":2", zone.Zonehash)

	zone, _ = findChatZone(root, "c23nb")
	zone.Count = maxRoomSize
	assert.Equal(t, ":5", zone.Zonehash)

	zone, _ = findChatZone(root, "c23nb")
	zone.Count = maxRoomSize
	assert.Equal(t, ":10", zone.Zonehash)

	zone, _ = findChatZone(root, "c23nb")
	zone.Count = maxRoomSize
	assert.Equal(t, ":21", zone.Zonehash)

	zone, _ = findChatZone(root, "c23nb")
	zone.Count = maxRoomSize
	assert.Equal(t, "c:1", zone.Zonehash)
}

func TestFindChatZone_SubhashToFirstGeohash2(t *testing.T) {
	root := createZone("", 1, nil)
	root.Count = maxRoomSize
	assert.Equal(t, ":1", root.Zonehash)

	zone, _ := findChatZone(root, "0000")
	zone.Count = maxRoomSize
	assert.Equal(t, ":2", zone.Zonehash)

	zone, _ = findChatZone(root, "0000")
	zone.Count = maxRoomSize
	assert.Equal(t, ":4", zone.Zonehash)

	zone, _ = findChatZone(root, "0000")
	zone.Count = maxRoomSize
	assert.Equal(t, ":8", zone.Zonehash)

	zone, _ = findChatZone(root, "0000")
	zone.Count = maxRoomSize
	assert.Equal(t, ":16", zone.Zonehash)

	zone, _ = findChatZone(root, "0000")
	zone.Count = maxRoomSize
	assert.Equal(t, "1:1", zone.Zonehash)
}

// func TestFindChatZone_PicksLeftChildFromGeohash(t *testing.T) {
// 	world := &Zone{Geohash: "", Count: 40}
// 	zone, _ := findChatZone(world, "c23nb")
// 	assert.Equal(t, world.Left, zone)
// }

// func TestFindChatZone_PicksRightChildFromGeohash(t *testing.T) {
// 	world := &Zone{Geohash: "", Count: 40}
// 	zone, _ := findChatZone(world, "k23nb")
// 	assert.Equal(t, world.Right, zone)
// }

// func TestFindChatZone_PicksLeftChildFromSubhash(t *testing.T) {
// 	world := &Zone{Geohash: "", Subhash: 0, Count: 30}
// 	zone, _ := findChatZone(world, "c23nb")
// 	assert.Equal(t, "", zone.Geohash)
// 	assert.Equal(t, 1, zone.Subhash)
// 	assert.Equal(t, world, zone.Parent)
// }

// func TestFindChatZone_SetsCorrectGeohashForChildren(t *testing.T) {
// 	root := &Zone{Geohash: "c23n", Subhash: 20, Count: 30}
// 	zone, _ := findChatZone(root, "c23nb")
// 	assert.Equal(t, "c23nb", zone.Geohash)
// }

// func TestFindChatZone_OverpopulatedRoom(t *testing.T) {
// 	root := &Zone{Geohash: "c23n", Subhash: 20, Count: 60}
// 	zone, err := findChatZone(root, "c23n")
// 	assert.Equal(t, root, zone)
// 	assert.Error(t, err, "Room full")
// }
