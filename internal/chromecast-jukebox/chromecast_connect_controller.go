package chromecastjukebox

import (
	"github.com/pkg/errors"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

var (
	connectNamespace = "urn:x-cast:com.google.cast.tp.connection"
)

type ChromecastConnectController struct {
	fromChromecastReader castchannel.Reader
	toChromecastWriter   castchannel.Writer
}

func NewChromecastConnectController() *ChromecastConnectController {
	fromChromecastReader, toChromecastWriter := castchannel.Pipe()

	return NewChromecastConnectControllerWithReaderWriter(fromChromecastReader, toChromecastWriter)
}

func NewChromecastConnectControllerWithReaderWriter(fromChromecastReader castchannel.Reader, toChromecastWriter castchannel.Writer) *ChromecastConnectController {
	return &ChromecastConnectController{
		fromChromecastReader: fromChromecastReader,
		toChromecastWriter:   toChromecastWriter,
	}
}

func (t *ChromecastConnectController) Read(cm *castchannel.CastMessage) error {
	return t.fromChromecastReader.Read(cm)
}

func (t *ChromecastConnectController) Write(cm *castchannel.CastMessage) error {
	return nil
}

func (t *ChromecastConnectController) Connect() error {
	payload := `{"type": "CONNECT"}`

	cm := castchannel.CastMessage{
		ProtocolVersion: castchannel.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &defaultSource,
		DestinationId:   &defaultDestination,
		Namespace:       &pingNamespace,
		PayloadType:     castchannel.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payload,
	}

	if err := t.toChromecastWriter.Write(&cm); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func (t *ChromecastConnectController) Close() error {
	payload := `{"type": "CLOSE"}`

	cm := castchannel.CastMessage{
		ProtocolVersion: castchannel.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &defaultSource,
		DestinationId:   &defaultDestination,
		Namespace:       &pingNamespace,
		PayloadType:     castchannel.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payload,
	}

	if err := t.toChromecastWriter.Write(&cm); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
