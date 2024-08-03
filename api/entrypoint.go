package api

import (
	"database/sql"
	"encoding/json"
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

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	fmt.Println("CORS middleware")
	return func(c *gin.Context) {
		// print the url of the request
		fmt.Printf("Request URL: %s\n", c.Request.URL)
		fmt.Printf("Request Method: %s\n", c.Request.Method)
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("Access-Control-Max-Age", "300")
	}
}

var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = 5432
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DATABASE")
	sslmode  = os.Getenv("POSTGRES_SSLMODE")
)

func init() {
	// err := godotenv.Load(".env")

    // if err != nil {
    //     log.Fatal("Error loading .env file")
    // }

	host = os.Getenv("POSTGRES_HOST")
	port = 5432
	user = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname = os.Getenv("POSTGRES_DATABASE")
	sslmode = os.Getenv("POSTGRES_SSLMODE")

	app = gin.New()
	r := app.Group("/api")
	r.Use(corsMiddleware())
	r.POST("/pageview", handlePageView)
	r.POST("/event", handleCustomEvent)
	r.GET("/pageviews", handlePageViews)
	r.GET("/custom-events", handleCustomEvents)
	r.OPTIONS("/pageview", handleOptions)
	r.OPTIONS("/event", handleOptions)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var err error
	if db == nil {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)

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

func handleOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Max-Age", "300")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
		customEvent, err := analyticsService.GetAllCustomEvents()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, customEvent)
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
    var event analytics.CustomEvent

    // Parse the JSON body
    if err := c.ShouldBindJSON(&event); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Serialize the Data field to JSON
    dataJSON, err := json.Marshal(event.Data)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize event data"})
        return
    }

    // Insert into database
    _, err = db.Exec("INSERT INTO custom_events (name, timestamp, data) VALUES ($1, $2, $3)",
        event.Name, event.Timestamp, string(dataJSON))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Event recorded successfully"})
}