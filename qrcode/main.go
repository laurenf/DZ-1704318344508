package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/skip2/go-qrcode"
)

func main() {
	sizes := []qrcode.RecoveryLevel{qrcode.Low, qrcode.Medium, qrcode.High, qrcode.Highest}
	for _, s := range sizes {
		err := qrcode.WriteColorFile("https://example.org?queryparam=bar", s, 256, color.Black, color.White, fmt.Sprintf("qr-%d.png", s))
		if err != nil {
			log.Fatal(err)
		}
	}
}
