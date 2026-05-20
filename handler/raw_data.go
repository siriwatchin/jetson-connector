package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/siriwatchin/jetson-connector/model"
)

type RawDataHandler struct {
	DB *gorm.DB
}

type rawDataRequest struct {
	Data dataPayload `json:"data"`
}

type dataPayload struct {
	Timestamp string  `json:"timestamp"`
	Type      string  `json:"type"`
	Color     string  `json:"color"`
	Brand     string  `json:"brand"`
	X         float32 `json:"x"`
	Y         float32 `json:"y"`
	Width     float32 `json:"width"`
	Height    float32 `json:"height"`
	CameraID  string  `json:"camera_id"`
	JetsonID  string  `json:"jetson_id"`
	TrackID   string  `json:"track_id"`
}

func (h *RawDataHandler) Create(c *gin.Context) {
	var req rawDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detectedAt, err := time.Parse(time.RFC3339, req.Data.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timestamp"})
		return
	}

	trackID, err := uuid.Parse(req.Data.TrackID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid track_id"})
		return
	}

	x, y, w, h2 := req.Data.X, req.Data.Y, req.Data.Width, req.Data.Height
	record := model.VehicleLog{
		DetectedAt:  detectedAt,
		TrackID:     trackID,
		CameraID:    req.Data.CameraID,
		VehicleType: strPtr(req.Data.Type),
		Color:       strPtr(req.Data.Color),
		Brand:       strPtr(req.Data.Brand),
		PositionX:   &x,
		PositionY:   &y,
		BboxWidth:   &w,
		BboxHeight:  &h2,
	}

	if err := h.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
