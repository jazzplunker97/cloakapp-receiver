package handlers

import (
	"context"
	"net/http"
	"time"

	"cloakapp-receiver/db"
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type TelemetryHandler struct {
	DB *db.InfluxDB
}

func NewTelemetryHandler(database *db.InfluxDB) *TelemetryHandler {
	return &TelemetryHandler{DB: database}
}

type TelemetryRequest struct {
	Host    string                 `json:"host" binding:"required"`
	Package string                 `json:"package" binding:"required"`
	FullURL string                 `json:"full_url" binding:"required"`
	Data    map[string]interface{} `json:"data" binding:"required"`
}

// @Summary Receive telemetry data
// @Description Receive telemetry data from frontend scripts and store in InfluxDB
// @Tags Telemetry
// @Accept json
// @Produce json
// @Param telemetry body TelemetryRequest true "Telemetry Data"
// @Success 200 {object} map[string]string
// @Router /telemetry [post]
func (h *TelemetryHandler) Receive(c *gin.Context) {
	var request TelemetryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Detect metadata internally
	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	// Tags: host, package, full_url, ip, and user_agent are primary tags
	tags := map[string]string{
		"host":       request.Host,
		"package":    request.Package,
		"full_url":   request.FullURL,
		"ip":         ip,
		"user_agent": ua,
	}

	// Fields: data contains the telemetry metrics
	fields := make(map[string]interface{})
	for k, v := range request.Data {
		fields[k] = v
	}

	// Default measurement name
	measurement := "telemetry"

	// Create point and write
	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	if err := h.DB.WriteAPI.WritePoint(context.Background(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
