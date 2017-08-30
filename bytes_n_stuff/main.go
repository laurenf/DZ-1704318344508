package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
)

// WriteFile writes a file of n random uint8, returning the filename
func WriteFile(n int) (string, error) {
	rand.Seed(1)
	b := make([]uint8, n)
	for i := 0; i < n; i++ {
		v := uint8(rand.Intn(256))
		fmt.Printf("wrote: %d\n", v)
		b[i] = v
	}
	bytes := make([]byte, n)
	for i, v := range b {
		bytes[i] = byte(v)
	}
	tmpfile, err := os.Create(fmt.Sprintf("/tmp/byte_length_%d.bin", n))
	if err != nil {
		return "", fmt.Errorf("error making temp file due to %s", err)
	}
	_, err = tmpfile.Write(bytes)
	if err != nil {
		return "", fmt.Errorf("error writing file due to %s", err)
	}
	return tmpfile.Name(), nil
}

// Read reads from reader into buffer of exactly length bytes
func Read(r io.Reader, buf []byte) error {
	if _, err := r.Read(buf); err != nil {
		return err
	}
	for _, b := range buf {
		v := uint8(b)
		fmt.Printf("read back: %d\n", v)
	}
	return nil
}

func main() {
	n := 1000
	file, err := WriteFile(n)
	if err != nil {
		fmt.Printf("couldn't write file due to %s", err)
		return
	}
	fmt.Printf("wrote %d bytes to %s\n", n, file)
	b := make([]byte, n)
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("couldnt open file %s due to %s", file, err)
		return
	}
	defer f.Close()
	err = Read(f, b)
	if err != nil {
		fmt.Printf("couldnt read from file due to %s", err)
		return
	}
}
