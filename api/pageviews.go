package handler

import (
	"database/sql"
	"encoding/json"
	"go-analytics/internal/analytics"
	"net/http"
	"time"
)

func PageViews(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=go_analytics sslmode=disable")
	if err != nil {
		panic(err)
	}

	analyticsService := analytics.New(db)

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

	db.Close()
}