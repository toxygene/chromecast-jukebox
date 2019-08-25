package chromecast_jukebox

import (
	cast_channel "chromecast_jukebox/internal/cast-channel"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMultiplex(t *testing.T) {
	toChromecastChannel := make(chan *cast_channel.CastMessage)
	fromControllerOneChannel := make(chan *cast_channel.CastMessage)
	fromControllerTwoChannel := make(chan *cast_channel.CastMessage)

	fromControllerChannels := []<-chan *cast_channel.CastMessage{fromControllerOneChannel, fromControllerTwoChannel}

	cm1 := cast_channel.CastMessage{}
	cm2 := cast_channel.CastMessage{}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err := Multiplex(toChromecastChannel, fromControllerChannels)
		assert.NoError(t, err)

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		fromControllerOneChannel <- &cm1
		fromControllerTwoChannel <- &cm2

		wg.Done()
	}()

	go func() {
		rcm1 := <- toChromecastChannel
		assert.Equal(t, &cm1, rcm1)

		rcm2 := <- toChromecastChannel
		assert.Equal(t, &cm2, rcm2)

		wg.Done()
	}()

	wg.Wait()
}
