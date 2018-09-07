package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

//
func benchmarkLoginReq(serverAddr string, n, c int32, isRan bool) (elapsed time.Duration) {

	// serverAddr := "http://localhost:8080/login"
	// n := int32(1000)
	// c := int32(2000)//concurrency
	// isRan := true
	readyGo := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(int(c))

	remaining := n

	var transport http.RoundTripper = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          int(c),
		MaxIdleConnsPerHost:   int(c),
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	cliRoutine := func(no int32) {
		for atomic.AddInt32(&remaining, -1) >= 0 {
			// continue

			data := url.Values{}

			var buffer bytes.Buffer
			buffer.WriteString("mhh")
			// rand
			// if isRan {
			// 	buffer.WriteString(strconv.Itoa(rand.Intn(1000000)))
			// } else {
			// 	buffer.WriteString("1")
			// }
			buffer.WriteString(strconv.Itoa(int(no)))
			username := buffer.String()

			data.Set("username", username)
			data.Set("password", "a123456")
			// data.Set("nickname", "newbot")
			// req, err := http.NewRequest("GET", serverAddr, bytes.NewBufferString(data.Encode()))
			req, err := http.NewRequest("POST", serverAddr, bytes.NewBufferString(data.Encode()))
			// req.AddCookie(&http.Cookie{Name: "username", Value: username, Expires: time.Now().Add(120 * time.Second), Path: "/"})
			// req.AddCookie(&http.Cookie{Name: "token", Value: "test", Expires: time.Now().Add(120 * time.Second), Path: "/"})

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it workparam=value
			if err != nil {
				log.Println(err)
			}
			<-readyGo
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
			}
			_, err1 := ioutil.ReadAll(resp.Body)
			if err1 != nil {
				log.Println(err1)
			}
			defer resp.Body.Close()
		}

		wg.Done()
	}

	for i := int32(0); i < c; i++ {
		go cliRoutine(i)
	}

	close(readyGo)
	start := time.Now()
	wg.Wait()
	// log.Println(time.Since(start))
	return time.Since(start)
}

func main() {
	// benchmarkLoginReq()
	//over 100 crash
	num := 2
	concurrency := 5
	elapsed := benchmarkLoginReq("http://localhost:8080/login", int32(num), int32(concurrency), true)
	fmt.Printf("\tTotal Requests(%v) - Concurrency(%v) - Cost(%s) - QPS(%v/sec)\n",
		num, concurrency, elapsed, math.Ceil(float64(1000)/(float64(elapsed)/1000000000)))

}
