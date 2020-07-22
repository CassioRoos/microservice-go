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

	c.l.Debug("Handle DELETE id: ", id)
	err := c.cr.DeleteCar(id)

	switch err {
	case nil:
		rw.WriteHeader(http.StatusNoContent)
	case data.ErrCarNotFound:
		{
			c.l.Error("[ERROR] id not found")
			rw.WriteHeader(http.StatusNotFound)
			data.ToJSON(&GenericError{http.StatusNotFound, err.Error()}, rw)
			return
		}
	default:
		c.l.Error("[ERROR] Error deleting record", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{http.StatusInternalServerError, err.Error()}, rw)
		return
	}

}
