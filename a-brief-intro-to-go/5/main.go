package main

import (
	"io"
	"net/http"
)

func main() {
	h := handler{Message: "HTTP server!"}
	http.ListenAndServe(":8080", h)
}

type handler struct {
	Message string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, h.Message)
}
