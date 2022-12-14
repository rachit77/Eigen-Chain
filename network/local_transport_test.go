package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestConnect(t *testing.T) {
// 	tra := NewLocalTransport("A")
// 	trb := NewLocalTransport("B")

// 	tra.Connect(trb)
// 	trb.Connect(tra)

// 	assert.Equal(t, tra.(*LocalTransport).peers[trb.Addr()], trb)
// 	assert.Equal(t, trb.(*LocalTransport).peers[tra.Addr()], tra)
// }

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)
	assert.Equal(t, 1, 1)

	msg := []byte("gm world!!")

	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))
	assert.Equal(t, <-trb.Consume(), RPC{
		From:    tra.Addr(),
		Payload: msg,
	})
}
