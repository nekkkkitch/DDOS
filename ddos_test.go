package ddos

import (
	"log"
	"net/http"
	"testing"
	"time"
)

func Test(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1", nil)
	if err != nil {
		t.Error(err)
	}
	ddoser, _ := New(2, true)
	ddoser.Start(req)
	time.Sleep(time.Second)
	ddoser.Stop()
	total, success := ddoser.GetStats()
	log.Printf("Stats: %v, %v", total, success)
}
