package types

type LatLng interface {
	PubSubSerializable
	BroadcastSerializable

	Lat() float64
	Lng() float64
	Geohash() string
}

type LatLngJSON struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (json *LatLngJSON) Type() PubSubDataType {
	return "latlng"
}
