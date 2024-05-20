package models

import (
	"log/slog"
	"slices"
	"time"

	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
)

func getCombinedStatus(statuses ...string) string {
	// Map all the statuses to expected cases
	for i, status := range statuses {
		statuses[i] = MapPaStatus(status)
	}

	// If any status is cancelled the end status is 'cancelled'
	if slices.Contains[[]string](statuses, internal.StatusCancelled) {
		return internal.StatusCancelled
	}

	// If any status is started the end status is 'started'
	if slices.Contains[[]string](statuses, internal.StatusInProgress) || slices.Contains[[]string](statuses, "") {
		return internal.StatusInProgress
	}

	// If any status is in_review the end status is 'in_review'
	if slices.Contains[[]string](statuses, internal.StatusInReview) || slices.Contains[[]string](statuses, "") {
		return internal.StatusInReview
	}
	// get the number of completed tasks
	count := 0
	for _, status := range statuses {
		if status != internal.StatusCompleted {
			count++
		}
	}

	switch count {
	case 0:
		return internal.StatusCompleted
	case len(statuses):
		return internal.StatusNotStarted
	default:
		return internal.StatusInProgress

	}
}

func getTaskCreatedAt(times ...time.Time) time.Time {
	createdAt := time.Time{}
	for _, time := range times {
		if !time.IsZero() && createdAt.Before(time) {
			createdAt = time
		}
	}

	return createdAt
}

func GetReasonCodeForTask(codes ...string) string {
	for _, code := range codes {
		if code != "" {
			// Return the first code that is not empty
			return code
		}
	}

	return ""
}

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
