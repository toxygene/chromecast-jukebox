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

	connection  castchannel.ReadWriter
	controllers []castchannel.ReadWriter
	logEntry    *logrus.Entry
}

func NewChromecast(connection io.ReadWriter, logEntry *logrus.Entry) *Chromecast {
	connectionController := NewChromecastConnectController()
	heartbeatController := NewChromecastPingController()

	return &Chromecast{
		ConnectionController: connectionController,
		HeartbeatController:  heartbeatController,
		connection:           castchannel.NewIoReadWriter(connection),
		controllers: []castchannel.ReadWriter{
			connectionController,
			heartbeatController,
		},
		logEntry: logEntry,
	}
}

func (t *Chromecast) RegisterController(c castchannel.ReadWriter) {
	t.controllers = append(t.controllers, c)
}

func (t *Chromecast) Run(ready *sync.WaitGroup) error {
	if err := t.ConnectionController.Connect(); err != nil {
		return errors.Wrap(err, "")
	}

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

		})
	}

	{
		for _, c := range t.controllers {
			g.Add(func() error {
				for {
					err := func(c castchannel.ReadWriter) error {
						var cm castchannel.CastMessage
						for {
							if err := c.Read(&cm); err != nil {
								return errors.Wrap(err, "")
							}

							if err := t.connection.Write(&cm); err != nil {
								return errors.Wrap(err, "")
							}
						}
					}(c)

					if err != nil {
						return err
					}
				}
			}, func(error) {

			})
		}
	}

	ready.Done()

	if err := g.Run(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
