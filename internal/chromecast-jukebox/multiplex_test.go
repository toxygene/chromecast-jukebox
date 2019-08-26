package chromecastjukebox

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toxygene/chromecast-jukebox/internal/cast-channel"
)

func TestMultiplex(t *testing.T) {
	toChromecastChannel := make(chan *castchannel.CastMessage)
	fromControllerOneChannel := make(chan *castchannel.CastMessage)
	fromControllerTwoChannel := make(chan *castchannel.CastMessage)

	fromControllerChannels := []<-chan *castchannel.CastMessage{fromControllerOneChannel, fromControllerTwoChannel}

	cm1 := castchannel.CastMessage{}
	cm2 := castchannel.CastMessage{}

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
		rcm1 := <-toChromecastChannel
		assert.Equal(t, &cm1, rcm1)

		rcm2 := <-toChromecastChannel
		assert.Equal(t, &cm2, rcm2)

		wg.Done()
	}()

	wg.Wait()
}
