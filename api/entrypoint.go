package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bartosz-skejcik/go-analytics/analytics"
	"github.com/bartosz-skejcik/go-analytics/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var analyticsService *analytics.Analytics
var app *gin.Engine
var db *sql.DB

var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = 5432
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DATABASE")
)

func init() {
	app = gin.New()
	r := app.Group("/api")
	r.POST("/pageview", handlePageView)
	r.POST("/event", handleCustomEvent)
	r.GET("/pageviews", handlePageViews)
	r.GET("/custom-events", handleCustomEvents)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var err error
	if db == nil {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=require",
			host, port, user, password, dbname)

		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			return
		}

		err = initializeTables(db)
		if err != nil {
			http.Error(w, "Failed to initialize tables", http.StatusInternalServerError)
			return
		}

		analyticsService = analytics.New(db)
	}

	app.ServeHTTP(w, r)
}

func initializeTables(db *sql.DB) error {
	query := `
	-- Create page_views table
	CREATE TABLE IF NOT EXISTS page_views (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL,
		timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
		user_agent TEXT,
		ip TEXT,
		referrer TEXT,
		country TEXT,
		country_code TEXT,
		os TEXT,
		browser TEXT
	);

	-- Create custom_events table
	CREATE TABLE IF NOT EXISTS custom_events (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
		data JSONB
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_page_views_timestamp ON page_views(timestamp);
	CREATE INDEX IF NOT EXISTS idx_custom_events_name_timestamp ON custom_events(name, timestamp);
	`

	_, err := db.Exec(query)
	return err
}


func handlePageViews(c *gin.Context) {
	startTime := time.Now().AddDate(0, 0, -1)
	endTime := time.Now()

	if c.Query("start") != "" {
		startTime, _ = time.Parse(time.RFC3339, c.Query("start"))
	}

	if c.Query("end") != "" {
		endTime, _ = time.Parse(time.RFC3339, c.Query("end"))
	}

	pageViews, err := analyticsService.GetPageViews(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pageViews)
}

func handleCustomEvents(c *gin.Context) {
	startTime := time.Now().AddDate(0, 0, -1)
	endTime := time.Now()
	name := c.Query("name")

	if c.Query("start") != "" {
		startTime, _ = time.Parse(time.RFC3339, c.Query("start"))
	}

	if c.Query("end") != "" {
		endTime, _ = time.Parse(time.RFC3339, c.Query("end"))
	}

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	customEvents, err := analyticsService.GetCustomEvents(name, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customEvents)
}

func handlePageView(c *gin.Context) {
	var pv analytics.PageView

	if err := c.ShouldBindJSON(&pv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ipInfo, err := utils.LocateIP(pv.IP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pv.Country = ipInfo.Country
	pv.CountryCode = ipInfo.CountryCode

	if err := analyticsService.RecordPageView(pv); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func handleCustomEvent(c *gin.Context) {
	var ce analytics.CustomEvent

	if err := c.ShouldBindJSON(&ce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := analyticsService.RecordCustomEvent(ce); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}