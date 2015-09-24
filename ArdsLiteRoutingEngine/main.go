// ArdsRoutingEngine project main.go
package main

import (
	"code.google.com/p/gorest"
	"fmt"
	"net/http"
	"time"
)

const layout = "2006-01-02T15:04:05Z07:00"

func main() {
	fmt.Println("Starting Ards Route Engine")
	InitiateRedis()
	go InitiateService()
	for {
		//fmt.Println("Searching...")
		availablePHashes := GetAllProcessingHashes()
		for _, h := range availablePHashes {
			if AcquireProcessingHashLock(h) == true {
				go ExecuteRequestHash(h)
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func InitiateService() {
	listeningPort := fmt.Sprintf(":%s", port)
	gorest.RegisterService(new(ArdsLiteRS))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(listeningPort, nil)
}
