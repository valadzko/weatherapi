package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/weather", weatherHandler)

	fmt.Println("Started server at port 8080:")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func weatherHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "test response for /weather")
}
