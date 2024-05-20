package models

import (
	"log/slog"

	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
)

func MapPaStatus(status string) string {
	switch status {
	case internal.StatusPending, internal.StatusNotStarted: // transaction:v2 event status
		return internal.StatusNotStarted
	case internal.StatusOpen, internal.StatusAccepted, internal.StatusInProgress: // transaction status or, transaction:v2 event status
		return internal.StatusInProgress
	case internal.StatusClosed, internal.StatusFulfilled, internal.StatusComplete, internal.StatusCompleted: // transaction is closed, event, expectation complete
		return internal.StatusCompleted
	case internal.StatusAborted, internal.StatusRejected, internal.StatusCancelled: // transaction is aborted, expectation skipped
		return internal.StatusCancelled
	case internal.StatusInReview:
		return internal.StatusInReview
	default:
		slogctx.Error(nil, "Unexpected status", slog.String("status", status))
	}

	return status
}

// Resolve any status conflicts with the consumer-api task status and what is coming from PA
func ResolveStatus(taskStatus, exptStatus string) string {
	switch taskStatus {
	case internal.StatusInReview, internal.StatusCancelled:
		return taskStatus

	default:
		return MapPaStatus(exptStatus)
	}
}
