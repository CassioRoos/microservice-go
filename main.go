package main

import (
	"MicroseService/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nicholasjackson/env"

	"github.com/gorilla/mux"
)

//A nice way to get the env variable, in this case, it will not raise an error when the value is not set, it will use default value instead
var bindAddress = env.String("BIND_ADDRESS", false, ":8888", "Bind address for the server")

func main() {
	env.Parse()
	log := log.New(os.Stdout, "cassio.roos-api++>", log.LstdFlags)
	car := handlers.NewCars(log)
	//Create a new serve mux and register the handler
	sm := mux.NewRouter()

	// SubRouter is a Handler of handler for GETs
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", car.GetCars)

	// SubRouter is a Handler of handler for PUTs
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	// Regex will be validated and the id value will be available in the service side
	putRouter.HandleFunc("/{id:[0-9]+}", car.UpdateCar)
	putRouter.Use(car.MiddlewareValidateCar)

	// SubRouter is a Handler of handler for POSTs
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", car.PostCar)
	postRouter.Use(car.MiddlewareValidateCar)

	// sm.Handle("/", car)

	server := &http.Server{
		Addr:         *bindAddress,
		Handler:      sm,
		ErrorLog:     log,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// log.Println("Listening to port: ", bindAddress)
	go func() {
		log.Printf("Starting server on port %s\n", *bindAddress)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	//creates a channel to listen OS signals in this case CTRL + C or when the program is killed
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Kill)
	signal.Notify(sigChan, os.Interrupt)

	// WAIT until the signal comes. This is blocking, then will wait until something occurs
	sig := <-sigChan
	log.Println("Shutdown gracefully", sig)

	// get the general context to create a new
	ct, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	server.Shutdown(ct)

}
