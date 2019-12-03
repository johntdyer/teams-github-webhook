package main

import (
	"fmt"
	resty "github.com/go-resty/resty"
	"log"
)

func dispatch(op, url, auth, body string) {
	client := resty.New()
	req := client.R()

	if body != "" {
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(body)
	}

	if auth != "" {
		req.SetHeader("Authorization", auth)
	}

        fmt.Println("URL:", url)
	if op == "POST" {
		resp, err := req.
			Post(url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Resp: ", resp)
	} else if op == "PUT" {
		resp, err := req.
			Put(url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Resp: ", resp)
	}
}
