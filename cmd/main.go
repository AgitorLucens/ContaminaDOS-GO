package main

import (
	"be/cmd/server"
	"context"
	"time"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Se crea contexto que escucha senal de interrupcion
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	//Se crea servidor en puerto establecido
	// HTTP only:
	// s := server.NewServer(":8080")
	// HTTP + HTTPS:
	s := server.NewServerWithTLS(":8080", ":8443", "cert.pem", "key.pem")
	//se corre servidor en un hilo propio
	go s.Run()

	<-ctx.Done()
	log.Println("\nReceived termination signal. Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Println("Error shutting down server:", err)
	} else {
		log.Println("Server shut down successfully.")
	}
}