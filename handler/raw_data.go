package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/siriwatchin/jetson-connector/model"
)

type RawDataHandler struct {
	DB          *gorm.DB
	EnableWrite bool
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
	raw, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))

	var req rawDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.saveFailedLog(raw, fmt.Sprintf("Create.ShouldBindJSON: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, err := toVehicleLog(req.Data)
	if err != nil {
		h.saveFailedLog(raw, fmt.Sprintf("Create.toVehicleLog: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.EnableWrite {
		if err := h.DB.Create(&record).Error; err != nil {
			h.saveFailedLog(raw, fmt.Sprintf("Create.DB.Create: %s", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type batchRequest struct {
	Data []dataPayload `json:"data"`
}

func (h *RawDataHandler) CreateBatch(c *gin.Context) {
	raw, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))

	var req batchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.saveFailedLog(raw, fmt.Sprintf("CreateBatch.ShouldBindJSON: %s", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	records := make([]model.VehicleLog, 0, len(req.Data))
	for _, payload := range req.Data {
		record, err := toVehicleLog(payload)
		if err != nil {
			h.saveFailedLog(raw, fmt.Sprintf("CreateBatch.toVehicleLog: %s", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		records = append(records, record)
	}

	if h.EnableWrite {
		if err := h.DB.CreateInBatches(records, 100).Error; err != nil {
			h.saveFailedLog(raw, fmt.Sprintf("CreateBatch.DB.CreateInBatches: %s", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *RawDataHandler) saveFailedLog(payload json.RawMessage, reason string) {
	if err := h.DB.Create(&model.FailedLog{Payload: payload, Reason: reason}).Error; err != nil {
		slog.Error("failed to insert failed_log", "error", err)
	}
}

func toVehicleLog(p dataPayload) (model.VehicleLog, error) {
	detectedAt, err := time.Parse(time.RFC3339, p.Timestamp)
	if err != nil {
		return model.VehicleLog{}, errors.New("invalid timestamp")
	}

	trackID := uuid.Nil
	if p.TrackID != "" {
		trackID, err = uuid.Parse(p.TrackID)
		if err != nil {
			return model.VehicleLog{}, errors.New("invalid track_id")
		}
	}

	x, y, w, h := p.X, p.Y, p.Width, p.Height
	return model.VehicleLog{
		DetectedAt:  detectedAt,
		TrackID:     trackID,
		CameraID:    p.CameraID,
		VehicleType: strPtr(p.Type),
		Color:       strPtr(p.Color),
		Brand:       strPtr(p.Brand),
		PositionX:   &x,
		PositionY:   &y,
		BboxWidth:   &w,
		BboxHeight:  &h,
	}, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
