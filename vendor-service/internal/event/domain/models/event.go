package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Outbox struct {
	Id          uuid.UUID       `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Exchange    string          `json:"exchange" gorm:"column:exchange;type:varchar(255);not null"`
	EventType   string          `json:"event_type" gorm:"column:event_type;type:varchar(255);not null"`
	Payload     json.RawMessage `json:"payload" gorm:"column:payload;type:jsonb;not null"`
	CreatedAt   time.Time       `gorm:"column:created_at;type:timestamp"`
	Processed   bool            `gorm:"column:processed;type:bool;not null;default:false"`
	ProcessedAt time.Time       `gorm:"column:processed_at;type:timestamp"`
}

type ProcessedMessage struct {
	MessageId uuid.UUID `json:"message_id" gorm:"column:message_id;type:uuid;primaryKey"`
}
