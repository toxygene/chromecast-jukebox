package chromecastjukebox

import (
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

// WriteCastMessage writes a Chromecast CastMessage to a writer
func WriteCastMessage(w io.Writer, message *castchannel.CastMessage) error {
	b, err := proto.Marshal(message)

	if err != nil {
		return errors.Wrap(err, "marshal message failed")
	}

	length := uint32(len(b))

	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return errors.Wrap(err, "write message length failed")
	}

	if err := binary.Write(w, binary.BigEndian, b); err != nil {
		return errors.Wrap(err, "write message bytes failed")
	}

	return nil
}

// ReadCastMessage reads a Chromecast CastMessage from a reader
func ReadCastMessage(r io.Reader) (*castchannel.CastMessage, error) {
	var expectedMessageLength uint32

	if err := binary.Read(r, binary.BigEndian, &expectedMessageLength); err != nil {
		return nil, errors.Wrap(err, "read expected message length from reader failed")
	}

	if expectedMessageLength > 0 {
		message := make([]byte, expectedMessageLength)

		messageLength, err := r.Read(message)

		if err != nil {
			return nil, errors.Wrap(err, "read message from reader failed")
		}

		if messageLength != int(expectedMessageLength) {
			return nil, errors.New("message length mismatch")
		}

		castMessage := &castchannel.CastMessage{}

		if err := proto.Unmarshal(message, castMessage); err != nil {
			return nil, errors.Wrap(err, "unmarshal message failed")
		}

		return castMessage, nil
	}

	return nil, nil
}

// GetCastMessagePayload unmarshals the payload of a CastMessage to a map[string]string
func GetCastMessagePayload(cm *castchannel.CastMessage, payload *map[string]string) error {
	if err := json.Unmarshal([]byte(*cm.PayloadUtf8), payload); err != nil {
		return errors.Wrap(err, "unmarshal payload failed")
	}

	return nil
}
