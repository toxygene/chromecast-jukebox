package castchannel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeWriteRead(t *testing.T) {
	r, w := Pipe()

	payload := `{"type": "LAUNCH", "requestId": 1, "appId": "1"}`
	sourceID := "source-1"
	destinationID := "receiver-1"
	namespace := "test"

	wcm := CastMessage{
		ProtocolVersion: CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &sourceID,
		DestinationId:   &destinationID,
		Namespace:       &namespace,
		PayloadType:     CastMessage_STRING.Enum(),
		PayloadUtf8:     &payload,
	}

	go func() {
		assert.NoError(t, w.Write(&wcm))
	}()

	var rcm CastMessage

	assert.NoError(t, r.Read(&rcm))
	assert.Equal(t, CastMessage_CASTV2_1_0.Enum(), rcm.ProtocolVersion)
	assert.Equal(t, &sourceID, rcm.SourceId)
	assert.Equal(t, &destinationID, rcm.DestinationId)
	assert.Equal(t, &namespace, rcm.Namespace)
	assert.Equal(t, CastMessage_STRING.Enum(), rcm.PayloadType)
	assert.Equal(t, &payload, rcm.PayloadUtf8)
}

func TestReadWrite(t *testing.T) {
	r, w := Pipe()

	payload := `{"type": "LAUNCH", "requestId": 1, "appId": "1"}`
	sourceID := "source-1"
	destinationID := "receiver-1"
	namespace := "test"

	go func() {
		var rcm CastMessage

		assert.NoError(t, r.Read(&rcm))
		assert.Equal(t, CastMessage_CASTV2_1_0.Enum(), rcm.ProtocolVersion)
		assert.Equal(t, &sourceID, rcm.SourceId)
		assert.Equal(t, &destinationID, rcm.DestinationId)
		assert.Equal(t, &namespace, rcm.Namespace)
		assert.Equal(t, CastMessage_STRING.Enum(), rcm.PayloadType)
		assert.Equal(t, &payload, rcm.PayloadUtf8)
	}()

	wcm := CastMessage{
		ProtocolVersion: CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &sourceID,
		DestinationId:   &destinationID,
		Namespace:       &namespace,
		PayloadType:     CastMessage_STRING.Enum(),
		PayloadUtf8:     &payload,
	}

	assert.NoError(t, w.Write(&wcm))
}
