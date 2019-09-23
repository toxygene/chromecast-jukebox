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

type Closer interface {
	Close() error
}

type ReadCloser interface {
	Reader
	Closer
}

type WriteCloser interface {
	Writer
	Closer
}

type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

type IoReadCloser struct {
	r io.ReadCloser
	l *logrus.Entry
}

func NewIoReadCloser(r io.ReadCloser, l *logrus.Entry) *IoReadCloser {
	return &IoReadCloser{r, l}
}

func (t *IoReadCloser) Read(cm *CastMessage) error {
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

type IoWriteCloser struct {
	w io.WriteCloser
	l *logrus.Entry
}

func NewIoWriteCloser(w io.WriteCloser, l *logrus.Entry) *IoWriteCloser {
	return &IoWriteCloser{w, l}
}

func (t *IoWriteCloser) Write(cm *CastMessage) error {
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

type IoReadWriteCloser struct {
	IoReadCloser
	IoWriteCloser
}

func NewIoReadWriteCloser(rw io.ReadWriteCloser, l *logrus.Entry) *IoReadWriteCloser {
	return &IoReadWriteCloser{
		IoReadCloser{rw, l},
		IoWriteCloser{rw, l},
	}
}

func (t *IoReadWriteCloser) Close() error {
	rErr := t.r.Close()
	wErr := t.w.Close()

	if rErr != nil || wErr != nil {
		return nil // todo
	}

	return nil
}
