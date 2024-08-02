package handler

import (
	"database/sql"
	"encoding/json"
	"go-analytics/internal/analytics"
	"go-analytics/internal/utils"
	"net/http"
)

func PageView(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=go_analytics sslmode=disable")
	if err != nil {
		panic(err)
	}

	analyticsService := analytics.New(db)

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

	db.Close()
}