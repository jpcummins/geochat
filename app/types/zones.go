package types

type Zones interface {
	Zone(string) (Zone, error)
	Refresh(string) (Zone, error)
	Save(Zone) error
}
