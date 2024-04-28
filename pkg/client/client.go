package client

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type ClientOptions struct {
	ServerAddr string
	Interval   int
}

func Run(serverAddr string, interval int) error {
	rand.Seed(time.Now().Unix())
	for range time.NewTicker(time.Duration(interval) * time.Second).C {
		url := fmt.Sprintf("%s/users/%d", serverAddr, rand.Intn(10))
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		bts, _ := io.ReadAll(resp.Body)
		log.Printf("Get %s, resp: %s", url, string(bts))
	}
	return nil
}
