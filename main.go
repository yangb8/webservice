package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yangb8/webservice/service"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/fib", service.FibHandler)
	//mux.HandleFunc("/_health", nil)
	//mux.HandleFunc("/_version", nil)
	//mux.HandleFunc("/", nil)

	srv := &http.Server{Addr: ":8080", Handler: mux}
	ssrv := &http.Server{Addr: ":8443", Handler: mux}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	go func() {
		<-sigCh
		// Stop catching signals, so that we can stop on second signal
		signal.Stop(sigCh)

		shutdown := func(s *http.Server) {
			defer wg.Done()
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			s.Shutdown(ctx)
		}

		log.Println("Shutting down ...")
		wg.Add(1)
		go shutdown(srv)
		wg.Add(1)
		go shutdown(ssrv)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Http Server: ", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := ssrv.ListenAndServeTLS("certs/server.crt", "certs/server.key"); err != nil {
			log.Println("Https Server: ", err)
		}
	}()
	wg.Wait()
}
