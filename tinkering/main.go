package main

import (
	"fmt"
	"log"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
)

func main() {

	kh, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DEBUG %#v\n", kh)
	fmt.Printf("DEBUG ks %#v\n", kh.ks)

	a, err := aead.New(kh)
	if err != nil {
		log.Fatal(err)
	}

	ct, err := a.Encrypt([]byte("this data needs to be encrypted"), []byte("associated data"))
	if err != nil {
		log.Fatal(err)
	}

	pt, err := a.Decrypt(ct, []byte("associated data"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cipher text: %s\nPlain text: %s\n", ct, pt)

}
