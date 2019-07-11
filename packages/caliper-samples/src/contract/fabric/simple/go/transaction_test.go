package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"testing"
)

func Must(err error) {
	panic(err)
}

func generateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	key, err := ecdsa.GenerateKey(S256(), rand.Reader)
	Must(err)

	return key, &key.PublicKey
}

func TestTransaction(t *testing.T) {
	//sk, vk := generateKeyPair()

}

func Example_parse_pem_with_btcec_public_uncompressed() {
	pubkey := []byte(`
-----BEGIN PUBLIC KEY-----
BPWvWexm6Ey/W26m16VWNjToErw8Dyp6hiWaGfsmd9l0KgtdRqwoHfvI3ZH+MlSKaewDtMlaZtLb9RsooaW+9p4=
-----END PUBLIC KEY-----`)

	pub, err := parse_pem_with_btcec_public_uncompressed(pubkey)
	fmt.Printf("Key: %+v, Err: %+v\n", pub, err)
	// Output:
}

// func BenchmarkSign(b *testing.B) {
// 	_, seckey := generateKeyPair()
// 	msg := csprngEntropy(32)
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		Sign(msg, seckey)
// 	}
// }
