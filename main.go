package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"pkg-viewer/controllers"
	"fmt"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.HomeHandler)

	fmt.Println("Listening on: 3000")

	http.ListenAndServe(":3000", r)
}