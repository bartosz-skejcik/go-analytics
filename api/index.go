package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/bartosz-skejcik/go-analytics/analytics"
	"github.com/bartosz-skejcik/go-analytics/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var analyticsService *analytics.Analytics

var (
	app *gin.Engine
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
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=go_analytics sslmode=disable")
	if err != nil {
		panic(err)
	}

	analyticsService = analytics.New(db)

	app.ServeHTTP(w, r)
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