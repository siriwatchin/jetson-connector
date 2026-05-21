package model

import (
	"encoding/json"
)

type FailedLog struct {
	Payload json.RawMessage `gorm:"column:payload;type:jsonb"`
	Reason  string          `gorm:"column:reason"`
}

func (FailedLog) TableName() string {
	return "failed_log"
}
