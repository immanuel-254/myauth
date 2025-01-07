package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/immanuel-254/myauth/internal"
)

func Api() {
	mux := http.NewServeMux()

	allviews := []internal.View{}

	internal.Routes(mux, allviews)

	server := &http.Server{
		Addr:         ":8080",                                                  // Custom port
		Handler:      internal.Cors(internal.New(internal.ConfigDefault)(mux)), // Attach the mux as the handler
		ReadTimeout:  10 * time.Second,                                         // Set read timeout
		WriteTimeout: 10 * time.Second,                                         // Set write timeout
		IdleTimeout:  30 * time.Second,                                         // Set idle timeout
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Error starting server:", err)
	}
}
