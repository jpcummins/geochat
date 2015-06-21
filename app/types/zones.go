package types

type Zones interface {
	Zone(string) (Zone, error)
	UpdateZone(string) (Zone, error)
	SetZone(Zone) error
}
