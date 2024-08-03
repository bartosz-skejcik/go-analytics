package main

import (
	"net/http"

	"github.com/bartosz-skejcik/go-analytics/api"
	"github.com/rs/cors"
)

func main() {

	c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000", "https://yourdomain.com"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
        // Debug: true, // Enable Debugging for testing, consider disabling in production
    })

	// apply the cors middleware
	http.Handle("/", c.Handler(http.HandlerFunc(api.Handler)))

	http.ListenAndServe(":8080", nil)
}