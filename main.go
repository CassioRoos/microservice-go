package main

import (
	"MicroseService/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	log := log.New(os.Stdout, "cassio.roos-api++>", log.LstdFlags)
	car := handlers.NewCars(log)

	sm := http.NewServeMux()
	sm.Handle("/", car)

	server := &http.Server{
		Addr:         ":8888",
		Handler:      sm,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  2 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Kill)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	log.Println("Shutdown gracefully", sig)
	ct, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ct)

}
