package chromecast_jukebox

import (
	"github.com/pkg/errors"
	cast_channel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
	"gopkg.in/tomb.v2"
)

func Distribute(fromChomecastChannel <-chan *cast_channel.CastMessage, toControllerChannels []chan<- *cast_channel.CastMessage) error {
	for cm := range fromChomecastChannel {
		err := func(cm *cast_channel.CastMessage) error {
			tb := tomb.Tomb{}

			for _, toControllerChannel := range toControllerChannels {
				func(toControllerChannel chan<- *cast_channel.CastMessage) {
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
