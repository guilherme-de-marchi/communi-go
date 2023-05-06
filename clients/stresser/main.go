package main

import (
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

func main() {
	start := time.Now()

	wg := new(sync.WaitGroup)
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go dnsRegisterStress(wg)
	}
	wg.Wait()

	slog.Info(
		"stress succedded",
		slog.Duration("time passed", time.Since(start)),
	)
}

func simpleStress(wg *sync.WaitGroup) {
	defer wg.Done()

	c, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		slog.Error(err.Error())
	}
	defer c.Close()
}

func dnsRegisterStress(wg *sync.WaitGroup) {
	defer wg.Done()

	c, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer c.Close()

	_, err = c.Write([]byte("set " + uuid.NewString() + " 123"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
