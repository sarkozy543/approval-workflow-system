package types

import (
	"encoding/json"
	"time"
)

// ✔ Status enum tipi
type ApprovalStatus string

const (
	StatusPending  ApprovalStatus = "PENDING"
	StatusApproved ApprovalStatus = "APPROVED"
	StatusRejected ApprovalStatus = "REJECTED"
)

// ✔ ApprovalRequest model
type ApprovalRequest struct {
	ID            string          `json:"id"`
	SourceEnv     string          `json:"source_env"`
	TargetEnv     string          `json:"target_env"`
	ChangePayload json.RawMessage `json:"change_payload"`
	Status        ApprovalStatus  `json:"status"`
	RequestedBy   string          `json:"requested_by"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
