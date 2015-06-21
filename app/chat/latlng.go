package chat

import (
	gh "github.com/TomiHiltunen/geohash-golang"
)

type LatLng struct {
	lat     float64
	lng     float64
	geohash string
}

func newLatLng(lat float64, lng float64) *LatLng {
	return &LatLng{
		lat:     lat,
		lng:     lng,
		geohash: gh.EncodeWithPrecision(lat, lng, 5),
	}
}

func (l *LatLng) Lat() float64 {
	return l.lat
}

func (l *LatLng) Lng() float64 {
	return l.lng
}

func (l *LatLng) Geohash() string {
	return l.geohash
}
