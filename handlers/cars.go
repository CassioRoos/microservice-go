package handlers

import (
	"MicroseService/data"
	"log"
	"net/http"
)

type (
	Cars struct {
		l *log.Logger
	}
)

func NewCars(l *log.Logger) *Cars {
	return &Cars{l}
}

func (c *Cars) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c.l.Println(r.Method)
	if r.Method == http.MethodGet {
		c.getCars(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		c.postCar(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (c *Cars) getCars(rw http.ResponseWriter, r *http.Request) {
	// Return the type CARS
	lc := data.GetCars()
	// FROM the type we can call the method to JSON, this way we can optimise our code
	if err := lc.ToJSON(rw); err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}
	// it's not necessary to call write, cuz we are already writing
	// rw.Write(data)
}

func (c *Cars) postCar(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Car POST")
	car := &data.Car{}
	err := car.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	c.l.Printf("Car %#v", car)
	data.AddCar(car)
}
