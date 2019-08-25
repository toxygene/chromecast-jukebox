package chromecast_jukebox

import (
	"io"

	cast_channel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/tomb.v2"
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

func (t *Chromecast) Run() error {
	tb := tomb.Tomb{}

	for _, c := range t.controllers {
		toChromecastChannel, _ := c.GetChannels()

		func(toChromecastChannel <-chan *cast_channel.CastMessage) {
			tb.Go(func() error {
				for {
					cm, ok := <-toChromecastChannel

					if !ok {
						return nil
					}

					// todo exclusive write lock needed here

					if err := WriteCastMessage(t.connection, cm); err != nil {
						return errors.Wrap(err, "")
					}
				}
			})
		}(toChromecastChannel)
	}

	tb.Go(func() error {
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
	})

	if err := tb.Wait(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
