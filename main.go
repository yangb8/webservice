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

	"github.com/yangb8/webservice/common/config"
	"github.com/yangb8/webservice/common/logger"
	"github.com/yangb8/webservice/common/sentry"
	"github.com/yangb8/webservice/common/statsd"
	"github.com/yangb8/webservice/service"
)

func main() {

	// load config
	cfg := config.GetConfig()

	// config log
	logger.SetLogger(cfg.Log.Location)

	// start sentry
	sentry.StartupSentry(cfg.Sentry.Enabled, cfg.Sentry.Dsn)
	defer sentry.ShutdownSentry()

	// start statsd
	statsd.StartupStatsd(cfg.Statsd.Enabled, cfg.Statsd.Address)
	defer statsd.ShutdownStatsd()

	// set routing
	rootMux := http.NewServeMux()
	rootMux.Handle("/compute/", http.StripPrefix("/compute", service.NewComputeHandler()))
	rootMux.Handle("/storage/", http.StripPrefix("/storage", service.NewStorageHandler(cfg)))
	rootMux.Handle("/_health", service.NewHealthHandler())

	// start services
	srv := &http.Server{Addr: ":8080", Handler: rootMux}
	ssrv := &http.Server{Addr: ":8443", Handler: rootMux}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		<-sigCh
		// Stop catching signals, so that we can stop on second signal
		signal.Stop(sigCh)

		shutdown := func(s *http.Server) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel() // make sure to cancel the context to avoid context leak
			s.Shutdown(ctx)
		}

		log.Println("Shutting down ...")
		go shutdown(srv)
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
