package pkg

import (
	"log"
	"sync"
	"testing"
)

func TestNewClient(t *testing.T) {
	var wg sync.WaitGroup
	count := 1000
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			client, _ := NewClient()
			log.Println(client.Id)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println("ok...")
}
