package chromecastjukebox

import (
	"io"
	"sync"

	"github.com/oklog/run"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	defaultSource      = "source-0"
	defaultDestination = "receiver-0"
)

type Chromecast struct {
	controllers []controller
	logEntry    *logrus.Entry
	connection  io.ReadWriteCloser
}

func NewChromecast(connection io.ReadWriteCloser, logEntry *logrus.Entry) *Chromecast {
	c := Chromecast{
		controllers: []controller{},
		logEntry:    logEntry,
		connection:  connection,
	}

	// register default controllers

	return &c
}

func (t *Chromecast) Close() error {
	if err := t.connection.Close(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func (t *Chromecast) Run() error {
	g := run.Group{}

	writeLock := sync.Mutex{}
	for _, c := range t.controllers {
		toChromecastChannel, _ := c.GetChannels()

		func(toChromecastChannel <-chan *castchannel.CastMessage) {
			g.Add(func() error {
				for {
					cm, ok := <-toChromecastChannel

					if !ok {
						return nil
					}

					writeLock.Lock()

					if err := WriteCastMessage(t.connection, cm); err != nil {
						writeLock.Unlock()
						return errors.Wrap(err, "")
					}

					writeLock.Unlock()
				}
			}, func(error) {
				c.Close()
			})
		}(toChromecastChannel)
	}

	g.Add(func() error {
		for {
			cm, err := ReadCastMessage(t.connection)

			if err != nil {
				return errors.Wrap(err, "")
			}

			for _, c := range t.controllers {
				func(c controller) {
					_, fromChromecastChannel := c.GetChannels()

					fromChromecastChannel <- cm
				}(c)
			}
		}
	}, func(error) {
		t.Close()
	})

	if err := g.Run(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
