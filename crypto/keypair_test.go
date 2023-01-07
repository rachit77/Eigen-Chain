package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestGenaratePrivateKey(t *testing.T) {
// 	privKey := GenaratePrivateKey()
// 	pubKey := privKey.PublicKey()
// 	address := pubKey.Address()

// 	msg := []byte("Hello World")
// 	sig, err := privKey.Sign(msg)
// 	assert.Nil(t, err)
// 	assert.True(t, sig.Verify(pubKey, msg))

// 	fmt.Println(address)
// 	fmt.Println(sig)
// }

func TestKeyPairSignVerifySuccess(t *testing.T) {
	privKey := GenaratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("Hello World")

	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	assert.True(t, sig.Verify(pubKey, msg))
}

func TestKeyPairSignVerifyFail(t *testing.T) {
	privKey := GenaratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("Hello World")

	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	OtherPrivKey := GenaratePrivateKey()
	OtherPublicKey := OtherPrivKey.PublicKey()
	assert.False(t, sig.Verify(OtherPublicKey, msg))
	assert.False(t, sig.Verify(pubKey, []byte("gm World!")))
}
