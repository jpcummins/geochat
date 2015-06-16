package types

type LatLng interface {
	Lat() float64
	Lng() float64
	Geohash() string
}
