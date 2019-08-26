package chromecastjukebox

import (
	"github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

type controller interface {
	GetChannels() (<-chan *castchannel.CastMessage, chan<- *castchannel.CastMessage)
	Run() error
}
