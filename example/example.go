package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/skaji/go-random-dialcontext"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: example URL")
		os.Exit(1)
	}
	url := os.Args[1]
	client := &http.Client{
		Transport: &http.Transport{
			DialContext:       random.DialContext(nil, nil),
			DisableKeepAlives: true,
		},
	}
	for i := 0; i < 10; i++ {
		res, err := client.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}
}
