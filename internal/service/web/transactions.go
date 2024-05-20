package web

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/juju/errors"
	errs "github.com/pkg/errors"
	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
)

func (s *WebService) GetTransaction(c *gin.Context, req *models.Transaction) (*models.TransactionResponse, error) {
	ctx := logTxReq(c, req.ID)

	tx, err := s.srv.GetTransaction(ctx, req)
	if err != nil {
		return nil, handleError(err)
	}

	return tx.ToResponse(), nil
}

func (s *WebService) GetTransactions(c *gin.Context) ([]models.TransactionResponse, error) {
	ctx := logTxReq(c, nil)

	tx, err := s.srv.GetTransactions(ctx)
	if err != nil {
		return nil, handleError(err)
	}
	return toTxResponses(tx, err)
}

func (s *WebService) PatchTransaction(c *gin.Context, req *models.TransactionStatus) (*models.TransactionResponse, error) {
	ctx := logTxReq(c, &req.ID)

	tx, err := s.srv.PatchTransaction(ctx, req)
	if err != nil {
		return nil, handleError(err)
	}

	return tx.ToResponse(), nil
}

func toTxResponses(res []models.Transaction, err error) ([]models.TransactionResponse, error) {
	err = handleError(err)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	var responses []models.TransactionResponse
	for _, tx := range res {
		t := tx.ToResponse()
		responses = append(responses, *t)
	}

	return responses, nil
}

func handleError(err error) error {
	if err != nil {
		if errs.Is(err, internal.ErrUnauthorized) {
			return errors.Unauthorized
		}
		if errs.Is(err, internal.ErrNotFound) {
			return errors.NotFoundf("Transaction")
		}
		if errs.Is(err, internal.ErrForbidden) {
			return errors.Forbidden
		}
		if errs.Is(err, internal.ErrBadRequest) {
			return errors.BadRequest
		}

		slogctx.Error(context.TODO(), "Unknown error in handleError", slog.Any("err", err))
		return errors.New("Internal server error")
	}

	return nil
}

// Add path variables to slogctx and headers to context
func logTxReq(c *gin.Context, txID *string) context.Context {
	ctx := getCtx(c)

	if txID != nil {
		ctx = slogctx.WithValue(ctx, "transactionID", txID)
	}

	headerMap := internal.GetContextHeaders(ctx)
	bearer := internal.SafeGetValueFromMap(headerMap, "Authorization")[len("Bearer "):]
	traceID := internal.SafeGetValueFromMap(headerMap, "Traceparent")
	if traceID == "" {
		slogctx.Warn(ctx, "Could not get traceID from headers",
			slog.Any("headers", headerMap),
		)
		traceID = uuid.New().String()
	}

	ctx = slogctx.WithValue(ctx, "traceId", traceID)

	// var uid string
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(bearer, &claims, nil)
	uid, ok := claims["user_id"].(string)
	if !ok {
		slogctx.Warn(ctx, "Could not get UID from token",
			slog.Any("claims", claims),
		)
	}

	if uid != "" {
		ctx = slogctx.WithValue(ctx, "UID", uid)
	}

	slogctx.Info(ctx, "Req",
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
	)

	return ctx
}

func getCtx(c *gin.Context) context.Context {
	ctx := c.Request.Context()

	slogctx.Debug(ctx, internal.HeaderKey, slog.Any(internal.HeaderKey, c.Request.Header))
	ctx = context.WithValue(ctx, internal.HeaderKey, c.Request.Header)

	return ctx
}
