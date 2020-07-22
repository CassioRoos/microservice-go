package handlers

import (
	"github.com/CassioRoos/MicroseService/data"
	"net/http"
)

// swagger:route GET /cars cars listCars
// Returns a list of cars
// responses:
// 		200: carsResponse

// ListAll handles GET requests and returns all current cars
func (c *Cars) GetListCars(rw http.ResponseWriter, r *http.Request) {
	c.l.Debug("Handle GET List Cars")
	rw.Header().Add("Content-Type", "application/json")
	cur := r.URL.Query().Get("currency")
	// Return the type CARS
	lc,err := c.cr.GetCars(cur)
	if err !=nil{
	    rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{http.StatusInternalServerError, err.Error()}, rw)
		return
	}
	// FROM the type we can call the method to JSON, this way we can optimise our code
	//r.Header.Add("Content-Type","application/json")

	if err := data.ToJSON(lc, rw); err != nil {
		c.l.Error("Unable to serialize car", "Error", err)
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}

	// it's not necessary to call write, cuz we are already writing
	// rw.Write(data)
}

func (c *Cars) GetCarById(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getCarId(r)
	c.l.Debug("Handle GET car by id:", id)
	cur := r.URL.Query().Get("currency")

	car, err := c.cr.GetCarById(id, cur)
	switch err {
	case nil:
	case data.ErrCarNotFound:
		{
			c.l.Error("[ERROR] fetching car", err)
			rw.WriteHeader(http.StatusNotFound)
			data.ToJSON(&GenericError{http.StatusNotFound, err.Error()}, rw)
			return
		}
	default:
		c.l.Error("[ERROR] fetching car", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{http.StatusInternalServerError, err.Error()}, rw)
		return
	}


	err = data.ToJSON(car, rw)
	if err != nil {
		// should never happen, but will log it - Defense coding
		c.l.Error("[ERROR] Serializing car ", err)
	}

}
