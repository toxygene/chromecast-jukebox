package chromecastjukebox

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

func TestConnect(t *testing.T) {
	r, w := castchannel.Pipe()

	c := NewChromecastConnectControllerWithReaderWriter(r, w)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		assert.NoError(t, c.Connect())
		wg.Done()
	}()

	var cm castchannel.CastMessage

	assert.NoError(t, c.Read(&cm))
	assert.Equal(t, `{"type": "CONNECT"}`, *cm.PayloadUtf8)

	wg.Wait()
}

func TestClose(t *testing.T) {
	r, w := castchannel.Pipe()

	c := NewChromecastConnectControllerWithReaderWriter(r, w)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		assert.NoError(t, c.Close())
		wg.Done()
	}()

	var cm castchannel.CastMessage

	assert.NoError(t, c.Read(&cm))
	assert.Equal(t, `{"type": "CLOSE"}`, *cm.PayloadUtf8)

	wg.Wait()
}

func TestRead(t *testing.T) {
	r, w := castchannel.Pipe()

	c := NewChromecastConnectControllerWithReaderWriter(r, w)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		assert.NoError(t, c.Connect())
		wg.Done()
	}()

	var cm castchannel.CastMessage

	assert.NoError(t, c.Read(&cm))
	assert.Equal(t, `{"type": "CONNECT"}`, *cm.PayloadUtf8)

	wg.Wait()
}
