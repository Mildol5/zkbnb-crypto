package twistedElgamal_bn128

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	"math/big"
	"testing"
)

func TestEncDec(t *testing.T) {
	sk, pk := GenKeyPair()
	fmt.Println("pk len:", len(pk.Bytes()))
	b := big.NewInt(int64(6))
	r, _ := rand.Int(rand.Reader, ORDER)
	enc := Enc(b, r, pk)
	encBytes, _ := json.Marshal(enc)
	fmt.Println("encBytes:", encBytes)
	var decodeEnc ElGamalEnc
	err := json.Unmarshal(encBytes, &decodeEnc)
	if err != nil {
		panic(err)
	}
	//dec := Dec(enc, sk)
	dec2 := DecByStart(&decodeEnc, sk, 0)
	//assert.Equal(t, b, dec)
	assert.Equal(t, b, dec2)
}
