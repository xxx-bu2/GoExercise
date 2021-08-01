package main

import (
	"context"
	"fmt"
	p_errors "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals)

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		select {
		case sig := <-signals:
			return fmt.Errorf("linux signals:%+v", sig)
		case <-ctx.Done():
			timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			if err := server.Shutdown(timeoutCtx); err != nil {
				return p_errors.Wrapf(ctx.Err(), "ctx error:%+v", ctx.Err())
			}
			return ctx.Err()
		}
	})

	if err := g.Wait(); err != nil {
		log.Printf("main exit:%+v", err)
	}
}
