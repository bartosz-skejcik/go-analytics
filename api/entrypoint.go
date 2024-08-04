// api/api.go
package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bartosz-skejcik/go-analytics/analytics"
	"github.com/bartosz-skejcik/go-analytics/utils"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var analyticsService *analytics.Analytics
var app *gin.Engine
var db *sql.DB

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
        c.Writer.Header().Set("Access-Control-Max-Age", "300")
    }
}

func init() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    app = gin.New()
    r := app.Group("/api")
    r.Use(corsMiddleware())
    r.POST("/session", handleSession)
    r.POST("/event", handleEvent)
    r.GET("/pageviews", handlePageViews)
    r.GET("/events", handleEvents)
    r.GET("/sessions", handleSessions)
    r.OPTIONS("/session", handleOptions)
    r.OPTIONS("/event", handleOptions)
}

func Handler(w http.ResponseWriter, r *http.Request) {
    var err error
    if db == nil {
        psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
            "password=%s dbname=%s sslmode=%s",
            os.Getenv("POSTGRES_HOST"), 5432, os.Getenv("POSTGRES_USER"),
            os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DATABASE"),
            os.Getenv("POSTGRES_SSLMODE"))

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
    CREATE TABLE IF NOT EXISTS sessions (
        id SERIAL PRIMARY KEY,
        anonymous_id TEXT NOT NULL,
        timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
        referrer TEXT,
        screen_width INTEGER,
        ip TEXT,
        user_agent TEXT,
        country TEXT,
        country_code TEXT,
        os TEXT,
        browser TEXT
    );

    CREATE TABLE IF NOT EXISTS page_views (
        id SERIAL PRIMARY KEY,
        session_id INTEGER NOT NULL REFERENCES sessions(id),
        url TEXT NOT NULL,
        view_order INTEGER NOT NULL,
        time_spent INTEGER NOT NULL
    );

    CREATE TABLE IF NOT EXISTS events (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
        data JSONB,
        ip TEXT,
        user_agent TEXT
    );

    CREATE INDEX IF NOT EXISTS idx_sessions_timestamp ON sessions(timestamp);
    CREATE INDEX IF NOT EXISTS idx_page_views_session_id ON page_views(session_id);
    CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
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

func handleSession(c *gin.Context) {
    var session analytics.Session
    if err := c.ShouldBindJSON(&session); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ipInfo, err := utils.LocateIP(session.IP)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    session.Country = ipInfo.Country
    session.CountryCode = ipInfo.CountryCode

    session.AnonymousID = utils.GenerateAnonymousID(session.IP, session.UserAgent)
    session.Timestamp = time.Now()

    if err := analyticsService.RecordSession(session); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusOK)
}

func handleSessions(c *gin.Context) {
    distinct := c.Query("distinct") == "true"
    sessions, err := analyticsService.GetSessions(distinct)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, sessions)
}

func handleEvent(c *gin.Context) {
    var event analytics.Event
    if err := c.ShouldBindJSON(&event); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    event.Timestamp = time.Now()

    if err := analyticsService.RecordEvent(event); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Event recorded successfully"})
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

func handleEvents(c *gin.Context) {
    startTime := time.Now().AddDate(0, 0, -1)
    endTime := time.Now()

    if c.Query("start") != "" {
        startTime, _ = time.Parse(time.RFC3339, c.Query("start"))
    }

    if c.Query("end") != "" {
        endTime, _ = time.Parse(time.RFC3339, c.Query("end"))
    }

    events, err := analyticsService.GetEvents(startTime, endTime)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, events)
}