package chromecastjukebox

import castchannel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"

type controller interface {
	Close() error
	GetChannels() (<-chan *castchannel.CastMessage, chan<- *castchannel.CastMessage)
	Run() error
}
