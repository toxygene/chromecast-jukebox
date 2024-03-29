package chromecastjukebox

import (
	"fmt"

	"github.com/pkg/errors"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

var (
	receiverNamespace = "urn:x-cast:com.google.cast.receiver"
)

type ChromecastReceiverController struct {
	toChromecastReader castchannel.ReadCloser
	toChromecastWriter castchannel.WriteCloser
}

// NewChromecastPingController creates a ChromecastPingController with a castchannel pipe for communication
func NewChromecastReceiverController() *ChromecastReceiverController {
	r, w := castchannel.Pipe()

	return NewChromecastReceiverControllerWithReaderWriter(r, w)
}

func NewChromecastReceiverControllerWithReaderWriter(r castchannel.ReadCloser, w castchannel.WriteCloser) *ChromecastReceiverController {
	return &ChromecastReceiverController{
		toChromecastReader: r,
		toChromecastWriter: w,
	}
}

func (t *ChromecastReceiverController) Close() error {
	rErr := t.toChromecastReader.Close()
	wErr := t.toChromecastWriter.Close()

	if rErr != nil || wErr != nil {
		return nil // todo
	}

	return nil
}

func (t *ChromecastReceiverController) Read(cm *castchannel.CastMessage) error {
	return t.toChromecastReader.Read(cm)
}

func (t *ChromecastReceiverController) Write(cm *castchannel.CastMessage) error {
	if cm.GetNamespace() != receiverNamespace {
		return nil
	}

	// var payload map[string]string

	// if err := json.Unmarshal([]byte(*cm.PayloadUtf8), &payload); err != nil {
	// 	return errors.Wrap(err, "")
	// }

	// todo

	return nil
}

func (t *ChromecastReceiverController) Launch(appID string) error {
	payload := fmt.Sprintf(
		`{"type": "LAUNCH", "requestId": 1, "appId": "%s"}`,
		appID,
	)

	cm := castchannel.CastMessage{
		ProtocolVersion: castchannel.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &defaultSource,
		DestinationId:   &defaultDestination,
		Namespace:       &receiverNamespace,
		PayloadType:     castchannel.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payload,
	}

	if err := t.toChromecastWriter.Write(&cm); err != nil {
		return errors.Wrap(err, "error writing launch to chromecast")
	}

	return nil
}
