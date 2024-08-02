package main

import (
	"database/sql"
	"encoding/json"
	"go-analytics/internal/analytics"
	"go-analytics/internal/utils"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var analyticsService *analytics.Analytics

func main() {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=go_analytics sslmode=disable")
	if err != nil {
		panic(err)
	}

	analyticsService = analytics.New(db)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/pageview", handlePageView)
	mux.HandleFunc("POST /api/event", handleCustomEvent)
	mux.HandleFunc("GET /api/pageviews", handlePageViews)
	mux.HandleFunc("GET /api/custom-events", handleCustomEvents)

	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Add your frontend URL here
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Wrap your mux with the CORS handler
	handler := c.Handler(mux)

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func handlePageViews(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().AddDate(0, 0, -1)
	endTime := time.Now()

	if r.URL.Query().Get("start") != "" {
		startTime, _ = time.Parse(time.RFC3339, r.URL.Query().Get("start"))
	}

	if r.URL.Query().Get("end") != "" {
		endTime, _ = time.Parse(time.RFC3339, r.URL.Query().Get("end"))
	}

	pageViews, err := analyticsService.GetPageViews(startTime, endTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pageViews)
}

func handleCustomEvents(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().AddDate(0, 0, -1)
	endTime := time.Now()
	name := ""

	if r.URL.Query().Get("start") != "" {
		startTime, _ = time.Parse(time.RFC3339, r.URL.Query().Get("start"))
	}

	if r.URL.Query().Get("end") != "" {
		endTime, _ = time.Parse(time.RFC3339, r.URL.Query().Get("end"))
	}

	if r.URL.Query().Get("name") != "" {
		name = r.URL.Query().Get("name")
	} else {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	customEvents, err := analyticsService.GetCustomEvents(name, startTime, endTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customEvents)
}

func handlePageView(w http.ResponseWriter, r *http.Request) {
	var pv analytics.PageView

	if err := json.NewDecoder(r.Body).Decode(&pv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ipInfo, err := utils.LocateIP(pv.IP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pv.Country = ipInfo.Country
	pv.CountryCode = ipInfo.CountryCode

	if err := analyticsService.RecordPageView(pv); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleCustomEvent(w http.ResponseWriter, r *http.Request) {
	var ce analytics.CustomEvent

	if err := json.NewDecoder(r.Body).Decode(&ce); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := analyticsService.RecordCustomEvent(ce); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}