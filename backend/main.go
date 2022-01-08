package main

import (
	"bakery-project/rest"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Staring REST User Service ...")
	a := rest.App{}
	a.Init("root", "admin123", "bakery")
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	a.Run(":8080")
	<-done
	log.Println("Stopping HTTP server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	if err := a.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v\n", err)
	}
	log.Print("Server exited properly")
}
