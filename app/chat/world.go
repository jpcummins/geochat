package chat

import (
	// "errors"
	// gh "github.com/TomiHiltunen/geohash-golang"
	// "github.com/jpcummins/geochat/app/cache"
	"github.com/jpcummins/geochat/app/types"
	// "strings"
)

var geohashmap = "0123456789bcdefghjkmnpqrstuvwxyz"

type World struct {
	root            types.Zone
	cache           types.Cache
	factory         types.Factory
	maxUsersPerZone int
	subscribe       <-chan types.Event
}

func newWorld(cache types.Cache, factory types.Factory, maxUsersPerZone int) (*World, error) {
	world := &World{
		cache:           cache,
		factory:         factory,
		maxUsersPerZone: maxUsersPerZone,
		subscribe:       make(<-chan types.Event),
	}

	root, err := world.GetOrCreateZone(":0z")
	if err != nil {
		return nil, err
	}

	world.root = root
	go world.manage()
	return world, nil
}

func (w *World) manage() { // It's a tough job.
	for {
		select {
		case event := <-w.subscribe:
			event.Data().OnReceive(event)
		}
	}
}

func (w *World) Factory() types.Factory {
	return w.factory
}

func (w *World) MaxUsersForNewZones() int {
	return w.maxUsersPerZone
}

func (w *World) GetOrCreateZone(id string) (types.Zone, error) {
	zone, err := w.cache.Zone(id)
	if err != nil {
		return nil, err
	}

	if zone == nil {
		zone, err = w.factory.NewZone(w, id)
		if err != nil {
			return nil, err
		}

		if err := w.cache.SetZone(zone); err != nil {
			return nil, err
		}
	}

	return zone, nil
}

//
// func getOrCreateAvailableZone(lat float64, long float64) (*Zone, error) {
// 	geohash := gh.EncodeWithPrecision(lat, long, 6)
// 	ch := make(chan interface{})
// 	world.getAvailableZone <- ch
// 	ch <- geohash
// 	zone := (<-ch).(*Zone)
// 	err := <-ch
// 	close(ch)
//
// 	if err != nil {
// 		return nil, err.(error)
// 	}
// 	return zone, nil
// }
//
// func GetOrCreateZone(id string) (*Zone, error) {
// 	ch := make(chan interface{})
// 	world.getZone <- ch
// 	ch <- id
// 	zone := (<-ch).(*Zone)
// 	err := <-ch
// 	close(ch)
//
// 	if err != nil {
// 		return nil, err.(error)
// 	}
//
// 	return zone, nil
// }
//
// func getOrCreateZone(id string) (*Zone, error) {
// 	// This algorithm is gross. I apologize if you have to read this.
// 	split := strings.Split(id, ":")
//
// 	// TODO: Validate string
// 	if len(split) != 2 || len(split[1]) != 2 {
// 		return nil, errors.New("Invalid id")
// 	}
//
// 	geohash := split[0]
// 	to := split[1][1]
//
// 	geohashLength := len(geohash)
// 	zone := world.root
//
// 	for {
// 		if zone.id == id {
// 			return zone, nil
// 		}
// 		zonegeohashLength := len(zone.geohash)
// 		if geohashLength > zonegeohashLength {
// 			if zone.left == nil || zone.right == nil {
// 				zone.createChildZones()
// 			}
//
// 			from_i := strings.Index(geohashmap, string(zone.from))
// 			to_i := strings.Index(geohashmap, string(zone.to))
// 			if to_i-from_i == 1 {
// 				if geohash[len(zone.geohash)] == zone.from {
// 					zone = zone.left
// 				} else {
// 					zone = zone.right
// 				}
// 				continue
// 			}
//
// 			if geohash[len(zone.geohash)] < zone.right.from {
// 				zone = zone.left
// 			} else {
// 				zone = zone.right
// 			}
// 		}
//
// 		if geohashLength == zonegeohashLength {
// 			if zone.left == nil || zone.right == nil {
// 				zone.createChildZones()
// 			}
// 			if to < zone.right.from {
// 				zone = zone.left
// 			} else {
// 				zone = zone.right
// 			}
// 		}
//
// 		if geohashLength < zonegeohashLength {
// 			return nil, errors.New("Error locating geohash: " + geohash)
// 		}
// 	}
// }
//
// func findChatZone(root *Zone, geohash string) (*Zone, error) {
// 	if root.left == nil && root.right == nil {
// 		root.createChildZones()
// 	}
//
// 	if root.isOpen {
// 		return root, nil
// 	}
//
// 	suffix := strings.TrimPrefix(geohash, root.geohash)
//
// 	if len(suffix) == 0 {
// 		return root, errors.New("Room full")
// 	}
//
// 	if root.geohash == root.right.geohash {
// 		if suffix[0] < root.right.from {
// 			return findChatZone(root.left, geohash)
// 		} else {
// 			return findChatZone(root.right, geohash)
// 		}
// 	} else {
// 		if suffix[0] < root.right.geohash[len(root.right.geohash)-1] {
// 			return findChatZone(root.left, geohash)
// 		} else {
// 			return findChatZone(root.right, geohash)
// 		}
// 	}
// }
