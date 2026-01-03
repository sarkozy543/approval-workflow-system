package types

import "time"

type ApprovalLog struct {
	ID        string    `json:"id"`
	RequestID string    `json:"request_id"`
	Action    string    `json:"action"`     // e.g. CREATED / APPROVED / REJECTED
	ActionBy  string    `json:"action_by"`  // e.g. "yusuf"
	Note      string    `json:"note"`       // optional
	CreatedAt time.Time `json:"created_at"`
}
const (
    ActionCreated  = "CREATED"
    ActionApproved = "APPROVED"
    ActionRejected = "REJECTED"
)
