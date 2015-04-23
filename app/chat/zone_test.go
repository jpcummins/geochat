package chat

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestFindChatZone_EmptyWorld(t *testing.T) {
    world := &Zone{Geohash: "", Subhash: 0}
    zone, _ := findChatZone(world, "c23nb")
    assert.Equal(t, world, zone)
    assert.Equal(t, "", world.Left.Geohash)
    assert.Equal(t, 1, world.Left.Subhash)
    assert.Equal(t, "", world.Right.Geohash)
    assert.Equal(t, 2, world.Right.Subhash)
}

func TestFindChatZone_PicksLeftChildFromGeohash(t *testing.T) {
    world := &Zone{Geohash: "", Count: 40}
    zone, _ := findChatZone(world, "c23nb")
    assert.Equal(t, world.Left, zone)
}

func TestFindChatZone_PicksRightChildFromGeohash(t *testing.T) {
    world := &Zone{Geohash: "", Count: 40}
    zone, _ := findChatZone(world, "k23nb")
    assert.Equal(t, world.Right, zone)
}

func TestFindChatZone_PicksLeftChildFromSubhash(t *testing.T) {
    world := &Zone{Geohash: "", Subhash: 0, Count: 30}
    zone, _ := findChatZone(world, "c23nb")
    assert.Equal(t, "", zone.Geohash)
    assert.Equal(t, 1, zone.Subhash)
    assert.Equal(t, world, zone.Parent)
}

func TestFindChatZone_SetsCorrectGeohashForChildren(t *testing.T) {
    root := &Zone{Geohash: "c23n", Subhash: 20, Count: 30}
    zone, _ := findChatZone(root, "c23nb")
    assert.Equal(t, "c23nb", zone.Geohash)
}

func TestFindChatZone_OverpopulatedRoom(t *testing.T) {
    root := &Zone{Geohash: "c23n", Subhash: 20, Count: 60}
    zone, err := findChatZone(root, "c23n")
    assert.Equal(t, root, zone)
    assert.Error(t, err, "Room full")
}
