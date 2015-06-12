package chat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	world := newWorld(nil)
	assert.Equal(t, ":0z", world.root.id)
	assert.Equal(t, 0, 0)
}
