package internal

import "time"

// PA/CE statuses
const (
	StatusOpen      = "open"
	StatusClosed    = "closed"
	StatusSubmitted = "submitted"
	StatusPending   = "pending"
	StatusComplete  = "complete"
	StatusCompleted = "completed" // used by apps
	StatusAccepted  = "accepted"
	StatusCancelled = "cancelled" // used by apps
	StatusFulfilled = "fulfilled"
	StatusError     = "error"
	StatusRejected  = "rejected"
	StatusAborted   = "aborted"
)

// statuses we use for FE
const (
	StatusNotStarted = "not_started" // pending
	StatusInProgress = "in_progress" // open
	StatusInReview   = "in_review"   // submitted
	// StatusCompleted  = "completed"
	// StatusCancelled  = "cancelled"
)

// date and time related variables
const (
	TimeYear      = 356 * 24 * time.Hour // Helper for calculating age
	DateLayout    = "2006-01-02"
	ActorCacheTTL = 15 * time.Minute
)

// Misc constants
const (
	HeaderKey = "headers"
)
