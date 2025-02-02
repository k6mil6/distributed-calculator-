package model

import (
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/internal/timeout"
	"time"
)

type Expression struct {
	ID         uuid.UUID       `json:"id"`
	UserID     int64           `json:"user_id"`
	Expression string          `json:"expression"`
	CreatedAt  time.Time       `json:"created_at"`
	Timeouts   timeout.Timeout `json:"timeouts"`
	IsTaken    bool            `json:"is_taken"`
	Result     float64         `json:"result,omitempty"`
	IsDone     bool            `json:"is_done,omitempty"`
}

type Subexpression struct {
	ID             int       `json:"id"`
	ExpressionId   uuid.UUID `json:"expression_id"`
	WorkerId       int       `json:"worker_id,omitempty"`
	Subexpression  string    `json:"subexpression"`
	IsTaken        bool      `json:"is_taken"`
	TakenAt        time.Time `json:"taken_at,omitempty"`
	IsBeingChecked bool      `json:"is_being_checked,omitempty"`
	Timeout        float64   `json:"timeout"`
	DependsOn      []int     `json:"depends_on,omitempty"`
	Result         float64   `json:"result,omitempty"`
	IsDone         bool      `json:"is_done,omitempty"`
}

type Timeouts struct {
	ID       int             `json:"id,omitempty"`
	UserID   int64           `json:"user_id,omitempty"`
	Timeouts timeout.Timeout `json:"timeouts,omitempty"`
}

type User struct {
	ID           int64  `json:"id"`
	Login        string `json:"login"`
	PasswordHash []byte `json:"password_hash"`
}

type Heartbeat struct {
	WorkerID int       `json:"worker_id"`
	SentAt   time.Time `json:"sent_at"`
}
