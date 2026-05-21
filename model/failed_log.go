package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type FailedLog struct {
	ID        uuid.UUID       `gorm:"column:id"`
	CreatedAt time.Time       `gorm:"column:created_at"`
	Payload   json.RawMessage `gorm:"column:payload;type:jsonb"`
	Reason    string          `gorm:"column:reason"`
}

func (FailedLog) TableName() string {
	return "failed_log"
}
