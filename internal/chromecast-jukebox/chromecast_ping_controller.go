package chromecast_jukebox

import (
	cast_channel "chromecast_jukebox/internal/cast-channel"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/tomb.v2"
)

var (
	pingNamespace = "urn:cast-cast:com.google.cast.tp.heartbeat"
)

type ChromecastPingController struct {
	pingTimer *time.Timer
	reader    chan *cast_channel.CastMessage
	writer    chan *cast_channel.CastMessage
}

func NewChromecastPingController(r chan *cast_channel.CastMessage, w chan *cast_channel.CastMessage) *ChromecastPingController {
	return &ChromecastPingController{
		pingTimer: time.NewTimer(5 * time.Second),
		reader:    r,
		writer:    w,
	}
}

func (t *ChromecastPingController) GetChannels() (<-chan *cast_channel.CastMessage, chan<- *cast_channel.CastMessage) {
	return t.reader, t.writer
}

func (t *ChromecastPingController) Close() error {
	close(t.reader)
	close(t.writer)

	return nil
}

func (t *ChromecastPingController) Run() error {
	tb := tomb.Tomb{}

	tb.Go(func() error {
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
	})

	tb.Go(func() error {
		for {
			<-t.pingTimer.C

			if err := t.Ping(); err != nil {
				return errors.Wrap(err, "")
			}
		}
	})

	if err := tb.Wait(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func (t *ChromecastPingController) Ping() error {
	payload := `{"type": "PING"}`

	cm := cast_channel.CastMessage{
		ProtocolVersion: cast_channel.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &defaultSource,
		DestinationId:   &defaultDestination,
		Namespace:       &pingNamespace,
		PayloadType:     cast_channel.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payload,
	}

	t.writer <- &cm

	return nil
}

func (t *ChromecastPingController) Pong() error {
	payload := `{"type": "PONG"}`

	cm := cast_channel.CastMessage{
		ProtocolVersion: cast_channel.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &defaultSource,
		DestinationId:   &defaultDestination,
		Namespace:       &pingNamespace,
		PayloadType:     cast_channel.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payload,
	}

	t.writer <- &cm

	return nil
}
