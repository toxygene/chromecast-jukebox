package chromecastjukebox

import (
	"sync"
	"testing"

	"github.com/toxygene/chromecast-jukebox/internal/cast-channel"

	"github.com/stretchr/testify/assert"
)

func TestDistribute(t *testing.T) {
	fromChromecastChannel := make(chan *castchannel.CastMessage)
	toControllerOneChannel := make(chan *castchannel.CastMessage)
	toControllerTwoChannel := make(chan *castchannel.CastMessage)

	toControllerChannels := []chan<- *castchannel.CastMessage{toControllerOneChannel, toControllerTwoChannel}

	cm := castchannel.CastMessage{}

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
