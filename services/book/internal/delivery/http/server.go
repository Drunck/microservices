package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type httpServer struct {
	router http.Handler
	config ServerConfig
}

type ServerConfig struct {
	Port         int
	IdleTimeout  string
	ReadTimeout  string
	WriteTimeout string
}

func NewHttpServer(router http.Handler, cfg ServerConfig) *httpServer {
	return &httpServer{router: router, config: cfg}
}

func (s *httpServer) Serve() error {

	readTimeout, err := time.ParseDuration(s.config.ReadTimeout)
	if err != nil {
		return err
	}

	writeTimeout, err := time.ParseDuration(s.config.WriteTimeout)
	if err != nil {
		return err
	}

	idleTimeout, err := time.ParseDuration(s.config.IdleTimeout)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", s.config.Port),
		Handler:      s.router,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Print("shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	fmt.Printf("starting server on %s", srv.Addr)

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Print("stopped server")

	return nil
}
