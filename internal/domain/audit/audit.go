package audit

import (
	"context"
	"time"
)

type Action string

const (
	ActionRotate    Action = "rotate"
	ActionView      Action = "view"
	ActionManualSet Action = "manual_set"
)

type Result string

const (
	ResultSuccess Result = "success"
	ResultFailure Result = "failure"
)

// AuditLog is a record of every operation Karden performs.
type AuditLog struct {
	ID        int
	TargetID  int
	Action    Action
	Actor     string
	Result    Result
	Reason    string
	CreatedAt time.Time
}

// Repository is the port for persisting AuditLogs.
type Repository interface {
	Save(ctx context.Context, log *AuditLog) error
	ListByTarget(ctx context.Context, targetID int) ([]*AuditLog, error)
}
