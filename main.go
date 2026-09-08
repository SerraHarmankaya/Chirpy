package main

import (
	"net/http"

	"github.com/SerraHarmankaya/Chirpy/system"
)

func main() {

	mux := http.ServeMux{}

	server := &http.Server{
		Addr:    ":8080",
		Handler: &mux,
	}

	fileServer := http.FileServer(http.Dir("./"))

	fileHandler := http.StripPrefix("/app", fileServer)

	mux.Handle("/app/", fileHandler)

	mux.HandleFunc("/healthz", system.ReadinessHandler)

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
