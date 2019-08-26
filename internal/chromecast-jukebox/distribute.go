package chromecastjukebox

import (
	"github.com/pkg/errors"
	"github.com/toxygene/chromecast-jukebox/internal/cast-channel"
	"gopkg.in/tomb.v2"
)

func Distribute(fromChomecastChannel <-chan *castchannel.CastMessage, toControllerChannels []chan<- *castchannel.CastMessage) error {
	for cm := range fromChomecastChannel {
		err := func(cm *castchannel.CastMessage) error {
			tb := tomb.Tomb{}

			for _, toControllerChannel := range toControllerChannels {
				func(toControllerChannel chan<- *castchannel.CastMessage) {
					tb.Go(func() error {
						toControllerChannel <- cm

						return nil
					})
				}(toControllerChannel)
			}

			if err := tb.Wait(); err != nil {
				return errors.Wrap(err, "")
			}

			return nil
		}(cm)

		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	return nil
}
