package domain

import "time"

type AuditLog struct {
	ID         int64
	EventType  string
	EntityType string
	EntityID   string
	Payload    []byte
	CreatedAt  time.Time
}
