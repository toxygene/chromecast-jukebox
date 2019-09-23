package chromecastjukebox

import (
	"io"
	"sync"

	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

var (
	defaultSource      = "source-0"
	defaultDestination = "receiver-0"
)

type Chromecast struct {
	ConnectionController *ChromecastConnectController
	HeartbeatController  *ChromecastPingController
	ReceiverController   *ChromecastReceiverController

	connection  castchannel.ReadWriteCloser
	controllers []castchannel.ReadWriteCloser
	l           *logrus.Entry
}

func NewChromecast(c io.ReadWriteCloser, l *logrus.Entry) *Chromecast {
	connectionController := NewChromecastConnectController()
	heartbeatController := NewChromecastPingController()
	receiverController := NewChromecastReceiverController()

	return &Chromecast{
		ConnectionController: connectionController,
		HeartbeatController:  heartbeatController,
		ReceiverController:   receiverController,
		connection:           castchannel.NewIoReadWriteCloser(c, l),
		controllers: []castchannel.ReadWriteCloser{
			connectionController,
			heartbeatController,
			receiverController,
		},
		l: l,
	}
}

func (t *Chromecast) Close() error {
	if err := t.ConnectionController.Disconnect(); err != nil {
		t.l.WithError(err).
			Error("error sending close command")

		return errors.Wrap(err, "error sending close command")
	}

	return nil
}

func (t *Chromecast) RegisterController(c castchannel.ReadWriteCloser) {
	t.controllers = append(t.controllers, c)
}

func (t *Chromecast) Run(ready *sync.WaitGroup) error {
	g := run.Group{}

	{
		g.Add(func() error {
			var cm castchannel.CastMessage
			for {
				if err := t.connection.Read(&cm); err != nil {
					return errors.Wrap(err, "error reading message from chromecast")
				}

				for _, c := range t.controllers {
					if err := c.Write(&cm); err != nil {
						return errors.Wrap(err, "error writing message to chromecast")
					}
				}
			}
		}, func(err error) {
			t.l.WithError(err).
				Error("interupting fan-out messages from chromecast")

			t.connection.Close()
		})
	}

	{
		mutex := sync.Mutex{}
		for _, c := range t.controllers {

			func(c castchannel.ReadWriteCloser) {
				g.Add(func() error {
					for {
						var cm castchannel.CastMessage

						if err := c.Read(&cm); err != nil {
							return errors.Wrap(err, "error reading message from controller")
						}

						mutex.Lock()
						if err := t.connection.Write(&cm); err != nil {
							mutex.Unlock()
							return errors.Wrap(err, "error writing message to chromecast")
						}
						mutex.Unlock()
					}
				}, func(err error) {
					t.l.WithField("controller", c).
						WithError(err).
						Error("intrupting fan-in messages from controller")

					c.Close()
				})
			}(c)
		}
	}

	// Send CONNECT
	{
		done := make(chan interface{})
		g.Add(func() error {
			if err := t.ConnectionController.Connect(); err != nil {
				t.l.WithError(err).
					Error("error connecting")

				return errors.Wrap(err, "error connecting")
			}

			ready.Done()

			<-done

			return nil
		}, func(err error) {
			t.l.WithError(err).
				Error("interupting connect controller")

			done <- struct{}{}
		})
	}

	if err := g.Run(); err != nil {
		return errors.Wrap(err, "error running chromecast")
	}

	return nil
}
