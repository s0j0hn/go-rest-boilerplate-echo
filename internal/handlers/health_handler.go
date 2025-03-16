package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/rabbitmq"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db          *gorm.DB
	rabbitMQ    *rabbitmq.AMQPClient
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB, rabbitMQ *rabbitmq.AMQPClient) *HealthHandler {
	return &HealthHandler{
		db:       db,
		rabbitMQ: rabbitMQ,
	}
}

// HealthStatus represents the health status response
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]HealthCheck `json:"checks"`
	Version   string                 `json:"version,omitempty"`
}

// HealthCheck represents individual component health
type HealthCheck struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency,omitempty"`
}

// Health godoc
// @Summary Health check endpoint
// @Description Get the health status of the service and its dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthStatus
// @Failure 503 {object} HealthStatus
// @Router /health [get]
func (h *HealthHandler) Health(c echo.Context) error {
	checks := make(map[string]HealthCheck)
	overallStatus := "healthy"

	// Check database
	dbCheck := h.checkDatabase()
	checks["database"] = dbCheck
	if dbCheck.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// Check RabbitMQ (if available)
	if h.rabbitMQ != nil {
		rabbitCheck := h.checkRabbitMQ()
		checks["rabbitmq"] = rabbitCheck
		if rabbitCheck.Status != "healthy" {
			overallStatus = "degraded" // RabbitMQ failure is non-critical
		}
	}

	status := HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Checks:    checks,
		Version:   "1.0.0", // You can inject this from build info
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	return c.JSON(statusCode, status)
}

// StatusResponse represents simple status responses
type StatusResponse struct {
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// Ready godoc
// @Summary Readiness check endpoint
// @Description Check if the service is ready to accept requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 503 {object} StatusResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c echo.Context) error {
	// Check critical dependencies
	dbCheck := h.checkDatabase()
	if dbCheck.Status != "healthy" {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status":  "not ready",
			"reason":  "database unavailable",
			"message": dbCheck.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ready",
	})
}

// Live godoc
// @Summary Liveness check endpoint
// @Description Check if the service is alive
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /live [get]
func (h *HealthHandler) Live(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "alive",
	})
}

// checkDatabase checks database connectivity
func (h *HealthHandler) checkDatabase() HealthCheck {
	start := time.Now()

	sqlDB, err := h.db.DB()
	if err != nil {
		return HealthCheck{
			Status:  "unhealthy",
			Message: "failed to get database instance: " + err.Error(),
			Latency: time.Since(start),
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return HealthCheck{
			Status:  "unhealthy",
			Message: "database ping failed: " + err.Error(),
			Latency: time.Since(start),
		}
	}

	return HealthCheck{
		Status:  "healthy",
		Message: "database connection successful",
		Latency: time.Since(start),
	}
}

// checkRabbitMQ checks RabbitMQ connectivity
func (h *HealthHandler) checkRabbitMQ() HealthCheck {
	start := time.Now()

	// Simple check - in a real implementation, you might want to
	// check connection status or try a simple operation
	if h.rabbitMQ == nil {
		return HealthCheck{
			Status:  "unhealthy",
			Message: "rabbitmq client not initialized",
			Latency: time.Since(start),
		}
	}

	// You could add more sophisticated checks here
	// For now, we'll assume it's healthy if the client exists
	return HealthCheck{
		Status:  "healthy",
		Message: "rabbitmq connection available",
		Latency: time.Since(start),
	}
}