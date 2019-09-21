package chromecastjukebox

import (
	"github.com/pkg/errors"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

var (
	connectNamespace = "urn:x-cast:com.google.cast.tp.connection"
)

type ChromecastConnectController struct {
	fromControllerReader castchannel.Reader
	toChromecastWriter   castchannel.Writer
}

func NewChromecastConnectController() *ChromecastConnectController {
	fromControllerReader, toChromecastWriter := castchannel.Pipe()

	return NewChromecastConnectControllerWithReaderWriter(fromControllerReader, toChromecastWriter)
}

func NewChromecastConnectControllerWithReaderWriter(fromControllerReader castchannel.Reader, toChromecastWriter castchannel.Writer) *ChromecastConnectController {
	return &ChromecastConnectController{
		fromControllerReader: fromControllerReader,
		toChromecastWriter:   toChromecastWriter,
	}
}

func (t *ChromecastConnectController) Read(cm *castchannel.CastMessage) error {
	return t.fromControllerReader.Read(cm)
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
		Namespace:       &connectNamespace,
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
		Namespace:       &connectNamespace,
		PayloadType:     castchannel.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payload,
	}

	if err := t.toChromecastWriter.Write(&cm); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
