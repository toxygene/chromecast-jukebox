package chromecastjukebox

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

func TestPingWrite(t *testing.T) {
	r, w := castchannel.Pipe()
	c := NewChromecastPingControllerWithReaderWriter(r, w)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		var rcm castchannel.CastMessage

		assert.NoError(t, c.Read(&rcm))

		payload := `{"type": "PING"}`

		assert.Equal(t, castchannel.CastMessage_CASTV2_1_0.Enum(), rcm.ProtocolVersion)
		assert.Equal(t, &pingNamespace, rcm.Namespace)
		assert.Equal(t, &defaultSource, rcm.SourceId)
		assert.Equal(t, &defaultDestination, rcm.DestinationId)
		assert.Equal(t, castchannel.CastMessage_STRING.Enum(), rcm.PayloadType)
		assert.Equal(t, &payload, rcm.PayloadUtf8)

		wg.Done()
	}()

	assert.NoError(t, c.Ping())

	wg.Wait()
}

func TestPongWrite(t *testing.T) {
	r, w := castchannel.Pipe()
	c := NewChromecastPingControllerWithReaderWriter(r, w)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		var rcm castchannel.CastMessage

		assert.NoError(t, c.Read(&rcm))

		payload := `{"type": "PONG"}`

		assert.Equal(t, castchannel.CastMessage_CASTV2_1_0.Enum(), rcm.ProtocolVersion)
		assert.Equal(t, &pingNamespace, rcm.Namespace)
		assert.Equal(t, &defaultSource, rcm.SourceId)
		assert.Equal(t, &defaultDestination, rcm.DestinationId)
		assert.Equal(t, castchannel.CastMessage_STRING.Enum(), rcm.PayloadType)
		assert.Equal(t, &payload, rcm.PayloadUtf8)

		wg.Done()
	}()

	assert.NoError(t, c.Pong())

	wg.Wait()
}

func TestPingRead(t *testing.T) {
	r, w := castchannel.Pipe()
	c := NewChromecastPingControllerWithReaderWriter(r, w)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		var rcm castchannel.CastMessage

		assert.NoError(t, r.Read(&rcm))

		payload := `{"type": "PONG"}`
		assert.Equal(t, &payload, rcm.PayloadUtf8)

		wg.Done()
	}()

	payload := `{"type": "PING"}`
	wcm := castchannel.CastMessage{
		PayloadUtf8: &payload,
	}

	assert.NoError(t, c.Write(&wcm))

	wg.Wait()
}
