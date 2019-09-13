package castchannel

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type Reader interface {
	Read(*CastMessage) error
}

type Writer interface {
	Write(*CastMessage) error
}

type ReadWriter interface {
	Reader
	Writer
}

type IoReader struct {
	r io.Reader
}

func NewIoReader(r io.Reader) *IoReader {
	return &IoReader{r}
}

func (t *IoReader) Read(cm *CastMessage) error {
	var expectedMessageLength uint32

	if err := binary.Read(t.r, binary.BigEndian, &expectedMessageLength); err != nil {
		return errors.Wrap(err, "read expected message length from reader failed")
	}

	if expectedMessageLength > 0 {
		message := make([]byte, expectedMessageLength)

		messageLength, err := t.r.Read(message)

		if err != nil {
			return errors.Wrap(err, "read message from reader failed")
		}

		if messageLength != int(expectedMessageLength) {
			return errors.New("message length mismatch")
		}

		if err := proto.Unmarshal(message, cm); err != nil {
			return errors.Wrap(err, "unmarshal message failed")
		}

		return nil
	}

	return nil
}

type IoWriter struct {
	w io.Writer
}

func NewIoWriter(w io.Writer) *IoWriter {
	return &IoWriter{w}
}

func (t *IoWriter) Write(cm *CastMessage) error {
	b, err := proto.Marshal(cm)

	if err != nil {
		return errors.Wrap(err, "marshal message failed")
	}

	length := uint32(len(b))

	if err := binary.Write(t.w, binary.BigEndian, length); err != nil {
		return errors.Wrap(err, "write message length failed")
	}

	if err := binary.Write(t.w, binary.BigEndian, b); err != nil {
		return errors.Wrap(err, "write message bytes failed")
	}

	return nil
}

type IoReadWriter struct {
	IoReader
	IoWriter
}

func NewIoReadWriter(rw io.ReadWriter) *IoReadWriter {
	return &IoReadWriter{
		IoReader{rw},
		IoWriter{rw},
	}
}
