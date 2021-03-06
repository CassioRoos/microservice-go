package main

import (
	"context"
	"fmt"
	"github.com/CassioRoos/MicroseService/data"
	"github.com/CassioRoos/MicroseService/grpc_healthcheck"
	"github.com/CassioRoos/MicroseService/handlers"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/nicholasjackson/env"

	protos "github.com/CassioRoos/grpc_currency/protos/currency"
	health "github.com/CassioRoos/grpc_currency/protos/healthcheck"
	gorilaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//A nice way to get the env variable, in this case, it will not raise an error when the value is not set, it will use default value instead
var bindAddress = env.String("APP_PORT", false, ":8888", "Bind address for the server")
var grpcPort = env.String("GRPC_PORT", false, "localhost:9098", "Bind address for GRPC server")

func main() {
	env.Parse()
	log := hclog.New(&hclog.LoggerOptions{
		Name:       "cassio.roos-api++>",
		Level:      hclog.LevelFromString("DEBUG"),
		JSONFormat: true,
		TimeFormat: "01/01/2006 15:04:05",
	})
	// THIS SHOULD NOT GO OUT IN PRODUCTION
	// FOR TESTING PURPOSES  ONLY
	log.Info("Establishing a connection to GRPC", "GRPC", *grpcPort)
	conn, err := grpc.Dial(*grpcPort, grpc.WithInsecure())
	if err != nil {
		log.Error("Unable to connect to GRPC Client on port")
		panic(err)
	}
	// I was having problems with GRPCurl in my containers
	// that`s why i create this work around
	hc := health.NewHealthCheckClient(conn)
	// simple struct with hclog
	healthCheck := grpc_healthcheck.NewGrpcHealthCheck(log, hc)
	// check 10 times for a OK status from GRPC, if the connection could not be established it will log an error and
	// will exit with code 1
	healthCheck.HealthCheck(10)
	cc := protos.NewCurrencyClient(conn)

	//log := log.New(os.Stdout, "cassio.roos-api++>", log.LstdFlags)
	validator := data.NewValidation()
	repo := data.NewCarsRepository(cc, log)
	car := handlers.NewCars(log, validator, repo)
	//Create a new serve mux and register the handler
	sm := mux.NewRouter()

	// SubRouter is a Handler of handler for GETs
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/cars", car.GetListCars)
	getRouter.HandleFunc("/cars", car.GetListCars).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/cars/{id:[0-9]+}", car.GetCarById)
	getRouter.HandleFunc("/cars/{id:[0-9]+}", car.GetCarById).Queries("currency", "{[A-Z]{3}}")

	// SubRouter is a Handler of handler for PUTs
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	// Regex will be validated and the id value will be available in the service side
	putRouter.HandleFunc("/cars", car.UpdateCar)
	putRouter.Use(car.MiddlewareValidateCar)

	// SubRouter is a Handler of handler for POSTs
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/cars", car.PostCar)
	postRouter.Use(car.MiddlewareValidateCar)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/cars/{id:[0-9]+}", car.DeleteCar)

	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)
	// Route to acces swager
	getRouter.Handle("/docs", sh)
	// Route to serve the yaml file to open-api
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// sm.Handle("/", car)
	// we could allow a specific host like http://localhost:3000
	ch := gorilaHandlers.CORS(gorilaHandlers.AllowedOrigins([]string{"*"}))
	server := &http.Server{
		Addr:         *bindAddress,
		Handler:      ch(sm),
		ErrorLog:     log.StandardLogger(&hclog.StandardLoggerOptions{}),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// log.Println("Listening to port: ", bindAddress)
	go func() {
		log.Info("Starting server on port %s\n", *bindAddress)
		if err := server.ListenAndServe(); err != nil {
			log.Error(fmt.Sprintf("Error while listening to port %s", *bindAddress))
			os.Exit(1)
		}
	}()

	//creates a channel to listen OS signals in this case CTRL + C or when the program is killed
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Kill)
	signal.Notify(sigChan, os.Interrupt)

	// WAIT until the signal comes. This is blocking, then will wait until something occurs
	sig := <-sigChan
	log.Info("Shutdown gracefully", sig)

	// get the general context to create a new
	ct, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	server.Shutdown(ct)

}
