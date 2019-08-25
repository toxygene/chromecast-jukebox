package chromecast_jukebox

import (
	cast_channel "chromecast_jukebox/internal/cast-channel"
	"github.com/pkg/errors"
	"gopkg.in/tomb.v2"
)

func Multiplex(toChromecastChannel chan<- *cast_channel.CastMessage, fromControllerChannels []<-chan *cast_channel.CastMessage) error {
	tb := tomb.Tomb{}

	for _, fromControllerChannel := range fromControllerChannels {
		func(fromControllerChannel <-chan *cast_channel.CastMessage) {
			tb.Go(func() error {
				cm, ok := <- fromControllerChannel
				if ok {
					toChromecastChannel <- cm
				}

				return nil
			})
		}(fromControllerChannel)
	}

	if err := tb.Wait(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}