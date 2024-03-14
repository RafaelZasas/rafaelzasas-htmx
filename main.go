package main

import (
	"context"
	"embed"
	"fmt"
	"htmx-go/core"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed views/*
var fs embed.FS

var app = core.NewApp(&fs)

func init() {
	if !app.IsBootstrapped() {
		if err := app.Bootstrap(); err != nil {
			log.Fatalf(fmt.Sprintf("init: failed to bootstrap app: %s", err))
		}
	}
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if !app.IsBootstrapped() {
		log.Fatalf("main: app is not bootstrapped")
	}

	if err := app.Start(ctx); err != nil {
		log.Fatalf(fmt.Sprintf("main: failed to start server: %s", err))
	}

	// Channel for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf(fmt.Sprintf("main: failed to shutdown app: %s", err))
	}
	log.Println("✅ App shutdown successful")
}
