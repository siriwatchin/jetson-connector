package model

import (
	"time"

	"github.com/google/uuid"
)

type VehicleLog struct {
	DetectedAt time.Time `gorm:"column:detected_at"`

	TrackID  uuid.UUID `gorm:"column:track_id"`
	CameraID string    `gorm:"column:camera_id"`

	LastSeen *time.Time `gorm:"column:last_seen"`

	VehicleType *string `gorm:"column:vehicle_type"`
	Color       *string `gorm:"column:color"`
	Brand       *string `gorm:"column:brand"`
	Plate       *string `gorm:"column:plate"`

	PositionX  *float32 `gorm:"column:position_x"`
	PositionY  *float32 `gorm:"column:position_y"`
	BboxWidth  *float32 `gorm:"column:bbox_width"`
	BboxHeight *float32 `gorm:"column:bbox_height"`

	EventType *string `gorm:"column:event_type"`

	TypeConfidence  *float32 `gorm:"column:type_confidence"`
	ColorConfidence *float32 `gorm:"column:color_confidence"`
	BrandConfidence *float32 `gorm:"column:brand_confidence"`
}

func (VehicleLog) TableName() string {
	return "vehicle_log"
}
