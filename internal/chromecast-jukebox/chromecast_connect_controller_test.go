package chromecastjukebox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

func TestConnect(t *testing.T) {
	r, w := castchannel.Pipe()

	c := NewChromecastConnectControllerWithReaderWriter(r, w)

	go func() {
		assert.NoError(t, c.Connect())
	}()

	var cm castchannel.CastMessage

	assert.NoError(t, c.Read(&cm))
	assert.Equal(t, `{"type": "CONNECT"}`, *cm.PayloadUtf8)
}

func TestClose(t *testing.T) {
	r, w := castchannel.Pipe()

	c := NewChromecastConnectControllerWithReaderWriter(r, w)

	go func() {
		assert.NoError(t, c.Close())
	}()

	var cm castchannel.CastMessage

	assert.NoError(t, c.Read(&cm))
	assert.Equal(t, `{"type": "CLOSE"}`, *cm.PayloadUtf8)
}

func TestRead(t *testing.T) {
	r, w := castchannel.Pipe()

	c := NewChromecastConnectControllerWithReaderWriter(r, w)

	go func() {
		assert.NoError(t, c.Connect())
	}()

	var cm castchannel.CastMessage

	assert.NoError(t, c.Read(&cm))
	assert.Equal(t, `{"type": "CONNECT"}`, *cm.PayloadUtf8)
}
