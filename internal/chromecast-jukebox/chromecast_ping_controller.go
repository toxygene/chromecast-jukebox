package chromecastjukebox

import (
	"time"

	"github.com/oklog/run"

	"github.com/toxygene/chromecast-jukebox/internal/cast-channel"

	"github.com/pkg/errors"
)

var (
	pingNamespace = "urn:cast-cast:com.google.cast.tp.heartbeat"
)

type ChromecastPingController struct {
	pingTimer *time.Timer
	reader    chan *castchannel.CastMessage
	writer    chan *castchannel.CastMessage
}

func NewChromecastPingController(r chan *castchannel.CastMessage, w chan *castchannel.CastMessage) *ChromecastPingController {
	return &ChromecastPingController{
		pingTimer: time.NewTimer(5 * time.Second),
		reader:    r,
		writer:    w,
	}
}

func (t *ChromecastPingController) GetChannels() (<-chan *castchannel.CastMessage, chan<- *castchannel.CastMessage) {
	return t.reader, t.writer
}

func (t *ChromecastPingController) Close() error {
	close(t.reader)
	close(t.writer)

	return nil
}

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
