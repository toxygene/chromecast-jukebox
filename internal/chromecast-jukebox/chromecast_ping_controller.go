package chromecastjukebox

import (
	"encoding/json"

	"github.com/pkg/errors"

	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

var (
	pingNamespace = "urn:x-cast:com.google.cast.tp.heartbeat"
)

// ChromecastPingController handles the ping/pong operations for a Chromecast communication session
type ChromecastPingController struct {
	toChromecastReader castchannel.ReadCloser
	toChromecastWriter castchannel.WriteCloser
}

// NewChromecastPingController creates a ChromecastPingController with a castchannel pipe for communication
func NewChromecastPingController() *ChromecastPingController {
	r, w := castchannel.Pipe()

	return NewChromecastPingControllerWithReaderWriter(r, w)
}

// NewChromecastPingControllerWithReaderWriter creates a ChromecastPingController using the supplied castchannel reader and writer for communication
func NewChromecastPingControllerWithReaderWriter(r castchannel.ReadCloser, w castchannel.WriteCloser) *ChromecastPingController {
	return &ChromecastPingController{
		toChromecastReader: r,
		toChromecastWriter: w,
	}
}

func (t *ChromecastPingController) Close() error {
	rErr := t.toChromecastReader.Close()
	wErr := t.toChromecastWriter.Close()

	if rErr != nil || wErr != nil {
		return nil // todo
	}

	return nil
}

func (t *ChromecastPingController) Read(cm *castchannel.CastMessage) error {
	return t.toChromecastReader.Read(cm)
}

func (t *ChromecastPingController) Write(cm *castchannel.CastMessage) error {
	if cm.GetNamespace() != pingNamespace {
		return nil
	}

	var payload map[string]string

	if err := json.Unmarshal([]byte(*cm.PayloadUtf8), &payload); err != nil {
		return errors.Wrap(err, "error unmarshaling message")
	}

	if payload["type"] == "PING" {
		t.Pong()
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

	if err := t.toChromecastWriter.Write(&cm); err != nil {
		return errors.Wrap(err, "error writing ping message to chromecast")
	}

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

	if err := t.toChromecastWriter.Write(&cm); err != nil {
		return errors.Wrap(err, "error writing pong message to chromecast")
	}

	return nil
}
