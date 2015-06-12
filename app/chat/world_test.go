package chat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWorld(t *testing.T) {
	world := newWorld(nil)
	assert.Equal(t, ":0z", world.root.id)
}
