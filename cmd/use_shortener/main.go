package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

func worker(i int, ch <-chan interface{}) {
	fmt.Printf("[%v]: Ready\n", i)
	<-ch
	fmt.Printf("[%v]: Started\n", i)
	bs := bytes.NewBufferString("http://ayaginkdkzmu.net/keu3mjdqmlun/jucsjdybso6s0")
	resp, err := http.Post("http://localhost:8080/", "text/plain; charset=utf-8", bs)
	if err != nil {
		panic(3)
	}
	defer resp.Body.Close()

	var rs bytes.Buffer
	rs.ReadFrom(resp.Body)
	fmt.Printf("[%v]: Done with %v\n", i, rs.String())
}

func main() {
	ch := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			worker(i, ch)
		}(i)
	}
	close(ch)
	wg.Wait()
}
