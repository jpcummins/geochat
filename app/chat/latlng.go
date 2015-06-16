package chat

import (
	gh "github.com/TomiHiltunen/geohash-golang"
)

type LatLng struct {
	lat     float64
	lng     float64
	geohash string
}

func newLatLng(latlng *gh.LatLng) *LatLng {
	return &LatLng{
		lat:     latlng.Lat(),
		lng:     latlng.Lng(),
		geohash: gh.EncodeWithPrecision(latlng.Lat(), latlng.Lng(), 5),
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
