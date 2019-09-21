package chromecastjukebox

import (
	"context"
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

	connection  castchannel.ReadWriter
	controllers []castchannel.ReadWriter
	l           *logrus.Entry
}

func NewChromecast(c io.ReadWriter, l *logrus.Entry) *Chromecast {
	connectionController := NewChromecastConnectController()
	heartbeatController := NewChromecastPingController()
	receiverController := NewChromecastReceiverController()

	return &Chromecast{
		ConnectionController: connectionController,
		HeartbeatController:  heartbeatController,
		ReceiverController:   receiverController,
		connection:           castchannel.NewIoReadWriter(c, l),
		controllers: []castchannel.ReadWriter{
			connectionController,
			heartbeatController,
			receiverController,
		},
		l: l,
	}
}

func (t *Chromecast) RegisterController(c castchannel.ReadWriter) {
	t.controllers = append(t.controllers, c)
}

func (t *Chromecast) Run(parentCtx context.Context, ready *sync.WaitGroup) error {
	g := run.Group{}

	{
		g.Add(func() error {
			var cm castchannel.CastMessage
			for {
				if err := t.connection.Read(&cm); err != nil {
					return errors.Wrap(err, "")
				}

				for _, c := range t.controllers {
					if err := c.Write(&cm); err != nil {
						return errors.Wrap(err, "")
					}
				}
			}
		}, func(error) {
			// todo - close the connection
		})
	}

	{
		for _, c := range t.controllers {
			func(c castchannel.ReadWriter) {
				g.Add(func() error {
					for {
						var cm castchannel.CastMessage

						if err := c.Read(&cm); err != nil {
							return errors.Wrap(err, "")
						}

						if err := t.connection.Write(&cm); err != nil {
							return errors.Wrap(err, "")
						}
					}
				}, func(error) {
					// todo - close the controller
				})
			}(c)
		}
	}

	// Send CONNECT
	{
		ctx, cancel := context.WithCancel(parentCtx)
		g.Add(func() error {
			if err := t.ConnectionController.Connect(); err != nil {
				return errors.Wrap(err, "")
			}

			ready.Done()

			<-ctx.Done()

			if err := t.ConnectionController.Close(); err != nil {
				return errors.Wrap(err, "")
			}

			return nil
		}, func(err error) {
			t.l.WithError(err).
				Error("interupting connect controller")

			cancel()
		})
	}

	// Send a PING every 5 seconds
	// {
	// 	ctx, cancel := context.WithCancel(parentCtx)
	// 	g.Add(func() error {
	// 		t.l.Trace("ping controller waiting for chromecast connection to be ready")

	// 		ready.Wait()

	// 		tick := time.NewTicker(5 * time.Second)
	// 		for {
	// 			select {
	// 			case <-ctx.Done():
	// 				return nil
	// 			case <-tick.C:
	// 				t.l.Trace("sending ping to chromecast")

	// 				t.HeartbeatController.Ping()
	// 			}
	// 		}
	// 	}, func(err error) {
	// 		t.l.WithError(err).
	// 			Error("intrupting ping controller timer")

	// 		cancel()
	// 	})
	// }

	if err := g.Run(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
