package main

import (
	"fmt"
	"os"
)

type message struct {
	res  chan *response
	req  chan *request
	quit chan int
}
type response struct {
	url string
	err interface{}
}
type request struct {
	url   string
	depth int
}

func newMessage() *message {
	return &message{
		res:  make(chan *response),
		req:  make(chan *request),
		quit: make(chan int),
	}
}

func (m *message) execute() {
	// ワーカーの数
	wc := 0
	urlMap := make(map[string]bool, 100)
	done := false
	for !done {
		select {
		case res := <-m.res:
			if res.err == nil {
				fmt.Printf("%s\n", res.url)
			} else {
				fmt.Fprintf(os.Stderr, "Error %s\n%v\n", res.url, res.err)
			}
		case req := <-m.req:
			if req.depth == 0 {
				break
			}

			if urlMap[req.url] {
				// 取得済み
				break
			}
			urlMap[req.url] = true

			wc++
			go Crawl(req.url, req.depth, m)
		case <-m.quit:
			wc--
			if wc == 0 {
				done = true
			}
		}
	}
	os.Exit(0)
}
