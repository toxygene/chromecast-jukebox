package chromecastjukebox

import (
	"github.com/pkg/errors"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

var (
	connectNamespace = "urn:x-cast:com.google.cast.tp.connection"
)

type ChromecastConnectController struct {
	fromControllerReader castchannel.ReadCloser
	toChromecastWriter   castchannel.WriteCloser
}

func NewChromecastConnectController() *ChromecastConnectController {
	fromControllerReader, toChromecastWriter := castchannel.Pipe()

	return NewChromecastConnectControllerWithReaderWriter(fromControllerReader, toChromecastWriter)
}

func NewChromecastConnectControllerWithReaderWriter(fromControllerReader castchannel.ReadCloser, toChromecastWriter castchannel.WriteCloser) *ChromecastConnectController {
	return &ChromecastConnectController{
		fromControllerReader: fromControllerReader,
		toChromecastWriter:   toChromecastWriter,
	}
}

func (t *ChromecastConnectController) Close() error {
	rErr := t.fromControllerReader.Close()
	wErr := t.toChromecastWriter.Close()

	if rErr != nil || wErr != nil {
		return nil // todo
	}

	return nil
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
		return errors.Wrap(err, "error writing connect message to chromecast")
	}

	return nil
}

func (t *ChromecastConnectController) Disconnect() error {
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
		return errors.Wrap(err, "error writing close message to chromecast")
	}

	return nil
}
