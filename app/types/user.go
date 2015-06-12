package types

import (
	"encoding/json"
)

type User interface {
	ID() string
	json.Marshaler
	json.Unmarshaler
}
