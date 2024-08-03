// analytics/analytics.go
package analytics

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Analytics struct {
    DB *sql.DB
}

type Session struct {
    AnonymousID string                 `json:"anonymous_id"`
    Timestamp   time.Time              `json:"timestamp"`
    Referrer    string                 `json:"referrer"`
    ScreenWidth int                    `json:"screen_width"`
    Pages       map[string][]int       `json:"pages"` // URL -> [order, time spent]
    IP          string                 `json:"ip"`
    UserAgent   string                 `json:"user_agent"`
    Country     string                 `json:"country"`
    CountryCode string                 `json:"country_code"`
    OS          string                 `json:"os"`
    Browser     string                 `json:"browser"`
}

type Event struct {
    Name      string                 `json:"name"`
    Timestamp time.Time              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
    IP        string                 `json:"ip"`
    UserAgent string                 `json:"user_agent"`
}

func New(db *sql.DB) *Analytics {
    return &Analytics{DB: db}
}

func (a *Analytics) RecordSession(s Session) error {
    tx, err := a.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Insert session data
    _, err = tx.Exec(`
        INSERT INTO sessions
        (anonymous_id, timestamp, referrer, screen_width, ip, user_agent, country, country_code, os, browser)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
        s.AnonymousID, s.Timestamp, s.Referrer, s.ScreenWidth, s.IP, s.UserAgent, s.Country, s.CountryCode, s.OS, s.Browser)
    if err != nil {
        return err
    }

    // Insert page views
    for url, data := range s.Pages {
        _, err = tx.Exec(`
            INSERT INTO page_views
            (anonymous_id, url, view_order, time_spent)
            VALUES ($1, $2, $3, $4)`,
            s.AnonymousID, url, data[0], data[1])
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (a *Analytics) RecordEvent(e Event) error {
    dataJSON, err := json.Marshal(e.Data)
    if err != nil {
        return err
    }

    _, err = a.DB.Exec(`
        INSERT INTO events
        (name, timestamp, data, ip, user_agent)
        VALUES ($1, $2, $3, $4, $5)`,
        e.Name, e.Timestamp, dataJSON, e.IP, e.UserAgent)
    return err
}

type PageViewResult struct {
    Timestamp   time.Time `json:"timestamp"`
    Referrer    string    `json:"referrer"`
    Country     string    `json:"country"`
    CountryCode string    `json:"country_code"`
    OS          string    `json:"os"`
    Browser     string    `json:"browser"`
    URL         string    `json:"url"`
    ViewOrder   int       `json:"view_order"`
    TimeSpent   int       `json:"time_spent"`
    IP          string    `json:"ip"`
    UserAgent   string    `json:"user_agent"`
}

func (a *Analytics) GetPageViews(startTime, endTime time.Time) ([]PageViewResult, error) {
    query := `
        SELECT s.timestamp, s.referrer, s.country, s.country_code, s.os, s.browser,
               pv.url, pv.view_order, pv.time_spent, s.ip, s.user_agent
        FROM sessions s
        JOIN page_views pv ON s.anonymous_id = pv.anonymous_id
        WHERE s.timestamp BETWEEN $1 AND $2
        ORDER BY s.timestamp DESC, pv.view_order
    `
    rows, err := a.DB.Query(query, startTime, endTime)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []PageViewResult
    for rows.Next() {
        var r PageViewResult
        err := rows.Scan(
            &r.Timestamp, &r.Referrer, &r.Country, &r.CountryCode,
            &r.OS, &r.Browser, &r.URL, &r.ViewOrder, &r.TimeSpent,
            &r.IP, &r.UserAgent,
        )
        if err != nil {
            return nil, err
        }
        results = append(results, r)
    }
    return results, nil
}

type EventResult struct {
    Name      string                 `json:"name"`
    Timestamp time.Time              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
    IP        string                 `json:"ip"`
    UserAgent string                 `json:"user_agent"`
}

func (a *Analytics) GetEvents(startTime, endTime time.Time) ([]EventResult, error) {
    query := `
        SELECT name, timestamp, data, ip, user_agent
        FROM events
        WHERE timestamp BETWEEN $1 AND $2
        ORDER BY timestamp DESC
    `
    rows, err := a.DB.Query(query, startTime, endTime)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []EventResult
    for rows.Next() {
        var r EventResult
        var dataJSON []byte
        err := rows.Scan(&r.Name, &r.Timestamp, &dataJSON, &r.IP, &r.UserAgent)
        if err != nil {
            return nil, err
        }
        err = json.Unmarshal(dataJSON, &r.Data)
        if err != nil {
            return nil, err
        }
        results = append(results, r)
    }
    return results, nil
}
