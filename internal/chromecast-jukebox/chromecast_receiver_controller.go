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
	toChromecastReader castchannel.Reader
	toChromecastWriter castchannel.Writer
}

// NewChromecastPingController creates a ChromecastPingController with a castchannel pipe for communication
func NewChromecastReceiverController() *ChromecastReceiverController {
	r, w := castchannel.Pipe()

	return NewChromecastReceiverControllerWithReaderWriter(r, w)
}

func NewChromecastReceiverControllerWithReaderWriter(r castchannel.Reader, w castchannel.Writer) *ChromecastReceiverController {
	return &ChromecastReceiverController{
		toChromecastReader: r,
		toChromecastWriter: w,
	}
}

func (t *ChromecastReceiverController) Read(cm *castchannel.CastMessage) error {
	return t.toChromecastReader.Read(cm)
}

func (t *ChromecastReceiverController) Write(cm *castchannel.CastMessage) error {
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
		return errors.Wrap(err, "")
	}

	return nil
}
