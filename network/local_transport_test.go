package network

import (
	"io/ioutil"
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
	msg := []byte("gm world!!")
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))

	rpc := <-trb.Consume()
	buf := make([]byte, len(msg))
	n, err := rpc.Payload.Read(buf)

	assert.Nil(t, err)
	assert.Equal(t, n, len(msg))

	assert.Equal(t, buf, msg)
	assert.Equal(t, rpc.From, tra.Addr())

}

func TestBroadast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")

	tra.Connect(trb)
	tra.Connect(trc)

	msg := []byte("gm! Eigen Chain")
	assert.Nil(t, tra.Broadcast(msg))

	rpcb := <-trb.Consume()
	b, err := ioutil.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, msg, b)

	rpcc := <-trc.Consume()
	c, err := ioutil.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, msg, c)
}
