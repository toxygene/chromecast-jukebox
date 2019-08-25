package chromecast_jukebox

import (
	cast_channel "github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

type controller interface {
	GetChannels() (<-chan *cast_channel.CastMessage, chan<- *cast_channel.CastMessage)
	Run() error
}
