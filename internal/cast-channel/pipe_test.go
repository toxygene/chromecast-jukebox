package castchannel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	r, w := Pipe()

	wcm := CastMessage{}

	go func() {
		assert.NoError(t, w.Write(&wcm))
	}()

	var rcm CastMessage

	assert.NoError(t, r.Read(&rcm))
	assert.Equal(t, wcm, rcm)
}
