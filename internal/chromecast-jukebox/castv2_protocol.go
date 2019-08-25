package chromecast_jukebox

import (
	cast_channel "chromecast_jukebox/internal/cast-channel"
	"encoding/binary"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"io"
)

func WriteCastMessage(w io.Writer, message *cast_channel.CastMessage) error {
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

func ReadCastMessage(r io.Reader) (*cast_channel.CastMessage, error) {
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

		castMessage := &cast_channel.CastMessage{}

		if err := proto.Unmarshal(message, castMessage); err != nil {
			return nil, errors.Wrap(err, "unmarshal message failed")
		}

		return castMessage, nil
	}

	return nil, nil
}

func GetCastMessagePayload(cm *cast_channel.CastMessage, payload *map[string]string) error {
	if err := json.Unmarshal([]byte(*cm.PayloadUtf8), payload); err != nil {
		return errors.Wrap(err, "unmarshal payload failed")
	}

	return nil
}