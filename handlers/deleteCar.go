package handlers

import (
	"github.com/CassioRoos/MicroseService/data"
	"net/http"
)

// swagger:route DELETE /car/{id} cars deleteCar
// Delete a car by the given Id
//
// responses:
//	201: noContentResponse
//	404: errorResponse
//	501: errorResponse
func (c *Cars) DeleteCar(rw http.ResponseWriter, r *http.Request) {
	id := getCarId(r)

	c.l.Println("Handle DELETE id: ", id)
	err := data.DeleteCar(id)

	switch err {
	case nil:
		rw.WriteHeader(http.StatusNoContent)
	case data.ErrCarNotFound:
		{
			c.l.Println("[ERROR] id not found")
			rw.WriteHeader(http.StatusNotFound)
			data.ToJSON(&GenericError{http.StatusNotFound, err.Error()}, rw)
			return
		}
	default:
		c.l.Println("[ERROR] Error deleting record", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{http.StatusInternalServerError, err.Error()}, rw)
		return
	}

}
