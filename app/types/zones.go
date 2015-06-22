package types

type Zones interface {
	FromCache(string) (Zone, error)
	FromDB(string) (Zone, error)
	Save(Zone) error
}
