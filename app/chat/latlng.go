package chat

import (
	gh "github.com/TomiHiltunen/geohash-golang"
	"github.com/jpcummins/geochat/app/types"
)

type LatLng struct {
	*types.LatLngJSON
	geohash string
}

func newLatLng(lat float64, lng float64) *LatLng {
	return &LatLng{
		LatLngJSON: &types.LatLngJSON{
			Lat: lat,
			Lng: lng,
		},
		geohash: gh.EncodeWithPrecision(lat, lng, 5),
	}
}

func (l *LatLng) Lat() float64 {
	return l.LatLngJSON.Lat
}

func (l *LatLng) Lng() float64 {
	return l.LatLngJSON.Lng
}

func (l *LatLng) Geohash() string {
	return l.geohash
}

func (l *LatLng) BroadcastJSON() interface{} {
	return l.LatLngJSON
}

func (l *LatLng) PubSubJSON() types.PubSubJSON {
	return l.LatLngJSON
}

func (l *LatLng) Update(types.PubSubJSON) error {
	return nil
}
