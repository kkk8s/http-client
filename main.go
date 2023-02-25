package main

import (
	"net/http"
	"os"
	"strconv"
	"time"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	timeout time.Duration
	url     string
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true,DisableQuote: true})
}

func main() {
	url = os.Getenv("REQUEST_URL")
	if url == "" {
		log.Fatalln("REQUEST_URL is not provided,program exit.")
	}
	if !strings.HasPrefix(url,"http://") {
		url = "http://" + url
	}

	// 超时处理，默认3s
	requestTimeout := os.Getenv("REQUEST_TIMEOUT")
	if requestTimeout == "" {
		requestTimeout = "3"
	}
	t, err := strconv.Atoi(requestTimeout)
	if err != nil {
		log.Fatalln("Parse timeout error,", err.Error())
	}
	timeout = time.Second * time.Duration(t)

	// 创建ticker对象
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	// 构造request及client对象
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Construct *http.Request error,", err.Error())
	}
	client := http.Client{Timeout: timeout}

	for {
		select {
		case <-ticker.C:
			_, err = client.Do(req)
			if err != nil {
				log.Printf("%v,timeout=%.1fs\n", err.Error(), client.Timeout.Seconds())
			}
		default:
		}
	}
}
