package chromecast_jukebox

import (
	cast_channel "chromecast_jukebox/internal/cast-channel"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistribute(t *testing.T) {
	fromChromecastChannel := make(chan *cast_channel.CastMessage)
	toControllerOneChannel := make(chan *cast_channel.CastMessage)
	toControllerTwoChannel := make(chan *cast_channel.CastMessage)

	toControllerChannels := []chan<- *cast_channel.CastMessage{toControllerOneChannel, toControllerTwoChannel}

	cm := cast_channel.CastMessage{}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		err := Distribute(fromChromecastChannel, toControllerChannels)
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rcm1 := <-toControllerOneChannel
		assert.Equal(t, &cm, rcm1)

		rcm2 := <-toControllerTwoChannel
		assert.Equal(t, &cm, rcm2)

		wg.Done()
	}()

	fromChromecastChannel <- &cm
	close(fromChromecastChannel)

	_, ok := <-fromChromecastChannel
	assert.False(t, ok)

	wg.Wait()
}
