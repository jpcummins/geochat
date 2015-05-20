package chat

// Zone information
// I'm not sure if this should be its own struct or if archive and subscribers
// should be members of Zone itself. Going with a seprate struct for now,
// because it is easiest to implement.
type ZoneInfo struct {
	ID          string         `json:"id"`
	Boundary    *ZoneBoundary  `json:"boundary"`
	Archive     *Archive       `json:"archive"`
	Subscribers []*Subscribers `json:"subscribers"`
}

func (z *ZoneInfo) Type() string {
	return "zone"
}
