package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHelloWold(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Something went wrong", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, "I received this (%s) in my body\n", data)
}
