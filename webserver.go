package main

import (
	"fmt"
	"log"
	"net/http"
)

func listPods(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Pods everywhere")
}

func main() {
	// New code
	http.HandleFunc("/", listPods)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
