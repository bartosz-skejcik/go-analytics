package analytics

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Analytics struct {
	DB *sql.DB
}

type PageView struct {
	URL       string `json:"url"`
	Timestamp string `json:"timestamp"`
	UserAgent string `json:"user_agent"`
	IP        string `json:"ip"`
	Referrer  string `json:"referrer"`
	Country   string `json:"country"`
	CountryCode string `json:"country_code"`
	OS        string `json:"os"`
	Browser   string `json:"browser"`
}

type CustomEvent struct {
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func New(db *sql.DB) *Analytics {
	return &Analytics{DB: db}
}

func (a *Analytics) RecordPageView(pv PageView) error {
	// Implementation to save page view to database
	if _, err := a.DB.Exec("INSERT INTO page_views (url, timestamp, user_agent, ip, referrer, country, country_code, os, browser) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", pv.URL, pv.Timestamp, pv.UserAgent, pv.IP, pv.Referrer, pv.Country, pv.CountryCode, pv.OS, pv.Browser); err != nil {
		return err
	} else {
		return nil
	}
}

func (a *Analytics) RecordCustomEvent(ce CustomEvent) error {
	// Implementation to save custom event to database
	if _, err := a.DB.Exec("INSERT INTO custom_events (name, timestamp, data) VALUES ($1, $2, $3)", ce.Name, ce.Timestamp, ce.Data); err != nil {
		return err
	} else {
		return nil
	}
}

func (a *Analytics) GetPageViews(startTime, endTime time.Time) ([]PageView, error) {
	query := `
		SELECT url, timestamp, user_agent, ip, referrer, country, os, browser, country_code
		FROM page_views
		ORDER BY timestamp DESC
	`
	rows, err := a.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pageViews []PageView
	for rows.Next() {
		var pv PageView
		err := rows.Scan(&pv.URL, &pv.Timestamp, &pv.UserAgent, &pv.IP, &pv.Referrer, &pv.Country, &pv.OS, &pv.Browser, &pv.CountryCode)
		if err != nil {
			return nil, err
		}
		pageViews = append(pageViews, pv)
	}

	return pageViews, nil
}

func (a *Analytics) GetCustomEvents(name string, startTime, endTime time.Time) ([]CustomEvent, error) {
	query := `
		SELECT name, timestamp, data
		FROM custom_events
		WHERE name = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp DESC
	`
	rows, err := a.DB.Query(query, name, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []CustomEvent
	for rows.Next() {
		var ce CustomEvent
		var dataJSON []byte
		err := rows.Scan(&ce.Name, &ce.Timestamp, &dataJSON)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(dataJSON, &ce.Data)
		events = append(events, ce)
	}

	return events, nil
}

func (a *Analytics) GetAllCustomEvents() ([]CustomEvent, error) {
	query := `
		SELECT name, timestamp, data
		FROM custom_events
		ORDER BY timestamp DESC
	`
	rows, err := a.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []CustomEvent
	for rows.Next() {
		var ce CustomEvent
		var dataJSON []byte
		err := rows.Scan(&ce.Name, &ce.Timestamp, &dataJSON)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(dataJSON, &ce.Data)
		events = append(events, ce)
	}

	return events, nil
}