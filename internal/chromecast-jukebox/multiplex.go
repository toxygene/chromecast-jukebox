package chromecastjukebox

import (
	"github.com/pkg/errors"
	castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
	"gopkg.in/tomb.v2"
)

func Multiplex(toChromecastChannel chan<- *castchannel.CastMessage, fromControllerChannels []<-chan *castchannel.CastMessage) error {
	tb := tomb.Tomb{}

	for _, fromControllerChannel := range fromControllerChannels {
		func(fromControllerChannel <-chan *castchannel.CastMessage) {
			tb.Go(func() error {
				cm, ok := <-fromControllerChannel
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
