package main

import (
	"context"
	"github.com/artushin/mailer/internal/http"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	log.Println("starting mailer service")
	ctx, cancel := context.WithCancel(context.Background())

	// start server
	http.Start(ctx, 8080)

	// watch for stop
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	sig := <-c

	switch sig {
	case syscall.SIGTERM:
		log.Println("server: catching sigterm")
		cancel()
		wait(time.Second * 10)
	case syscall.SIGINT:
		log.Println("server: catching sigint")
		cancel()
		wait(time.Second * 2)
	}
}

func wait(timeout time.Duration) {
	done := make(chan struct{}, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		<-http.Close()
		wg.Done()
	}()

	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		os.Exit(0)
	case <-time.After(timeout):
		log.Fatal("service timed out on shutdown")
	}
}
