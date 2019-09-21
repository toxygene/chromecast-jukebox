package castchannel

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	l *logrus.Entry
}

func NewIoReader(r io.Reader, l *logrus.Entry) *IoReader {
	return &IoReader{r, l}
}

func (t *IoReader) Read(cm *CastMessage) error {
	var expectedMessageLength uint32

	t.l.Info("reading expected message length from reader")

	if err := binary.Read(t.r, binary.BigEndian, &expectedMessageLength); err != nil {
		t.l.WithError(err).
			Error("read expected length from reader failed")

		return errors.Wrap(err, "read expected message length from reader failed")
	}

	t.l.WithField("expectedMessageLength", expectedMessageLength).
		Info("read expected message length")

	if expectedMessageLength > 0 {
		message := make([]byte, expectedMessageLength)

		t.l.Info("reading message from reader")

		messageLength, err := t.r.Read(message)
		if err != nil {
			t.l.WithError(err).
				Error("read message from reader failed")

			return errors.Wrap(err, "read message from reader failed")
		}

		if messageLength != int(expectedMessageLength) {
			t.l.Error("message length does not match expected message length")

			return errors.New("message length does not match expected message length")
		}

		l := t.l.WithField("message", string(message))

		l.Info("read message from reader")

		l.Info("unmarshaling protobuf message")

		if err := proto.Unmarshal(message, cm); err != nil {
			l.WithError(err).
				Error("unmarshal message failed")

			return errors.Wrap(err, "unmarshal message failed")
		}

		l.Info("read message from reader succeeded")

		return nil
	}

	return nil
}

type IoWriter struct {
	w io.Writer
	l *logrus.Entry
}

func NewIoWriter(w io.Writer, l *logrus.Entry) *IoWriter {
	return &IoWriter{w, l}
}

func (t *IoWriter) Write(cm *CastMessage) error {
	l := t.l.WithField("message", cm.String())

	l.Trace("marshaling protobuf message")

	b, err := proto.Marshal(cm)
	if err != nil {
		l.WithError(err).
			Error("marshaling protobuf message failed")

		return errors.Wrap(err, "marshaling protobuf message failed")
	}

	length := uint32(len(b))

	l.Trace("writing protobuf message length")

	if err := binary.Write(t.w, binary.BigEndian, length); err != nil {
		l.WithError(err).
			Error("write protobuf message length failed")

		return errors.Wrap(err, "write protobuf message length failed")
	}

	l.Trace("writing protobuf message")

	if err := binary.Write(t.w, binary.BigEndian, b); err != nil {
		l.WithError(err).
			Error("write protobuf message failed")

		return errors.Wrap(err, "write protobuf message failed")
	}

	l.Trace("protobuf message write succeeded")

	return nil
}

type IoReadWriter struct {
	IoReader
	IoWriter
}

func NewIoReadWriter(rw io.ReadWriter, l *logrus.Entry) *IoReadWriter {
	return &IoReadWriter{
		IoReader{rw, l},
		IoWriter{rw, l},
	}
}
