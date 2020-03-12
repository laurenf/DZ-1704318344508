package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"math/rand"
	"testing"
)

var key []byte
var plaintext []byte
var ciphertext []byte

func BenchmarkAES(b *testing.B) {
	plaintextSizes := []int{10, 100, 1000, 10_000, 100_000, 1_000_000, 10_000_000, 100_000_000}
	keySizes := []int{16, 32}
	// fixed nonce
	nonce := []byte("123456789123")
	for _, p := range plaintextSizes {
		for _, k := range keySizes {
			key = make([]byte, k)
			r := rand.New(rand.NewSource(1))
			_, err := r.Read(key)
			if err != nil {
				log.Fatalf("err getting random %s", err)
			}
			plaintext = make([]byte, p)
			_, err = r.Read(key)
			if err != nil {
				log.Fatalf("err getting random for plaintext %s", err)
			}
			block, err := aes.NewCipher(key)
			if err != nil {
				log.Fatalf("err creating cipher %s", err)
			}
			aesgcm, err := cipher.NewGCM(block)
			if err != nil {
				log.Fatalf("err creating ccm %s", err)
			}
			b.Run(fmt.Sprintf("plaintext size %d bytes, key size %d bytes", p, k), func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					ciphertext = aesgcm.Seal(nil, nonce, plaintext, nil)
				}
			})
		}
	}
}
