package domain

import "time"

// Action represents what operation was performed
type Action string

const (
	ActionRotate    Action = "rotate"
	ActionView      Action = "view"
	ActionManualSet Action = "manual_set"
)

// Result represents whether the operation succeeded
type Result string

const (
	ResultSuccess Result = "success"
	ResultFailure Result = "failure"
)

// AuditLog is a record of every operation Janusd performs
type AuditLog struct {
	ID        int
	TargetID  int
	Action    Action
	Actor     string
	Result    Result
	Reason    string // failure reason, empty on success
	RotatedAt time.Time
	Metadata  map[string]string
}
