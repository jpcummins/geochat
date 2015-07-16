package types

type Zones interface {
	Zone(string) (Zone, error)
	FromCache(string) Zone
	UpdateCache(Zone)
	FromDB(string) (Zone, error)
	Save(Zone) error
}
