package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New()

	resp, err := client.R().
		Get("https://httpbin.org/headers")

	if err != nil {
		fmt.Printf("error calling httpbin due to %s", err)
		return
	}
	fmt.Println(resp)
	fmt.Println(resp.Header())
}
