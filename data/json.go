package data

import (
	"encoding/json"
	"io"
)

func FromJSON(i interface{}, r io.Reader) error {
	// the opposite of encode
	d := json.NewDecoder(r)
	return d.Decode(i)
}

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	// now we have to encode ourselfs, because C is pointing to the slice of cars
	return e.Encode(i)
}
