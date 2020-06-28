package handlers

import (
	"MicroseService/data"
	"net/http"
)

// swagger:route GET /cars cars listCars
// Returns a list of cars
// responses:
// 		200: carsResponse

// ListAll handles GET requests and returns all current cars
func (c *Cars) GetListCars(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET List Cars")
	// Return the type CARS
	lc := data.GetCars()
	// FROM the type we can call the method to JSON, this way we can optimise our code
	//r.Header.Add("Content-Type","application/json")
	rw.Header().Set("Content-Type", "application/json")
	if err := data.ToJSON(lc, rw); err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}

	// it's not necessary to call write, cuz we are already writing
	// rw.Write(data)
}

func (c *Cars) GetCarById(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle GET car by id")
	id := getCarId(r)
	c.l.Println("Handle GET car by id:", id)

	car, err := data.GetCarById(id)
	switch err {
	case nil:
	case data.ErrCarNotFound:
		{
			c.l.Println("[ERROR] fetching car", err)
			rw.WriteHeader(http.StatusNotFound)
			data.ToJSON(&GenericError{http.StatusNotFound, err.Error()}, rw)
			return
		}
	default:
		c.l.Println("[ERROR] fetching car", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{http.StatusInternalServerError, err.Error()}, rw)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	err = data.ToJSON(car, rw)
	if err != nil {
		// should never happen, but will log it - Defense coding
		c.l.Println("[ERROR] Serializing car ", err)
	}

}
