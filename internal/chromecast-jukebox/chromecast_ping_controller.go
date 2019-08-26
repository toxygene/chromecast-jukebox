package chromecastjukebox

import (
	"time"

	"github.com/oklog/run"

	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"

	"github.com/pkg/errors"
)

var (
	pingNamespace = "urn:cast-cast:com.google.cast.tp.heartbeat"
)

// ChromecastPingController handles the ping/pong operations for a Chromecast communication session
type ChromecastPingController struct {
	pingTimer *time.Timer
	reader    chan *castchannel.CastMessage
	writer    chan *castchannel.CastMessage
}

// NewChromecastPingController constructs a ChromecastPingController, using the supplied channels for reading from and writing to the Chromecast
func NewChromecastPingController(r chan *castchannel.CastMessage, w chan *castchannel.CastMessage) *ChromecastPingController {
	return &ChromecastPingController{
		pingTimer: time.NewTimer(5 * time.Second),
		reader:    r,
		writer:    w,
	}
}

// GetChannels returns the controllers read from and write to Chromecast channels
func (t *ChromecastPingController) GetChannels() (<-chan *castchannel.CastMessage, chan<- *castchannel.CastMessage) {
	return t.reader, t.writer
}

// Close closes the read from and write to Chromecast channels
func (t *ChromecastPingController) Close() error {
	close(t.reader)
	close(t.writer)

	return nil
}

// Run causes the controller to listen for PING and PONG payloads and periodicly sends PING payloads to the Chromecast
func (t *ChromecastPingController) Run() error {
	g := run.Group{}

	g.Add(func() error {
		for {
			cm, ok := <-t.reader
			if !ok {
				return nil
			}

			var payload map[string]string

			if err := GetCastMessagePayload(cm, &payload); err != nil {
				return errors.Wrap(err, "")
			}

			if payload["type"] == "PING" {
				if err := t.Pong(); err != nil {
					return errors.Wrap(err, "")
				}
			}
		}
	}, func(error) {
		close(t.reader)
	})

	g.Add(func() error {
		for {
			<-t.pingTimer.C

			if err := t.Ping(); err != nil {
				return errors.Wrap(err, "")
			}
		}
	}, func(error) {
		t.pingTimer.Stop()
	})

	if err := g.Run(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

// Ping sends a PING payload to the Chromecast
func (t *ChromecastPingController) Ping() error {
	payload := `{"type": "PING"}`

	cm := castchannel.CastMessage{
		ProtocolVersion: castchannel.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &defaultSource,
		DestinationId:   &defaultDestination,
		Namespace:       &pingNamespace,
		PayloadType:     castchannel.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payload,
	}

	t.writer <- &cm

	return nil
}

// Pong sends a PONG payload to the Chromecast
func (t *ChromecastPingController) Pong() error {
	payload := `{"type": "PONG"}`

	cm := castchannel.CastMessage{
		ProtocolVersion: castchannel.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &defaultSource,
		DestinationId:   &defaultDestination,
		Namespace:       &pingNamespace,
		PayloadType:     castchannel.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payload,
	}

	t.writer <- &cm

	return nil
}
