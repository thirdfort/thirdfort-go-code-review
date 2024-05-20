package web

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	jujuerr "github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/pkg/errors"
	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
)

func ErrHook(c *gin.Context, e error) (int, interface{}) {
	slogctx.Debug(context.TODO(), "Handling errors", slog.Any("e", e))
	errcode, errpl := 500, e.Error()
	_, ok := e.(tonic.BindError)
	if ok {
		errcode, errpl = 400, e.Error()
	} else {
		switch {
		case errors.Is(e, internal.ErrBadRequest) || errors.Is(e, internal.ErrAlreadyExists) || errors.Is(e, jujuerr.BadRequest) || errors.Is(e, jujuerr.NotValid) || errors.Is(e, jujuerr.AlreadyExists) ||
			errors.Is(e, jujuerr.NotSupported) || errors.Is(e, jujuerr.NotAssigned) || errors.Is(e, jujuerr.NotProvisioned):
			errcode, errpl = 400, e.Error()
		case errors.Is(e, internal.ErrForbidden) || errors.Is(e, jujuerr.Forbidden):
			errcode, errpl = 403, e.Error()
		case errors.Is(e, jujuerr.MethodNotAllowed):
			errcode, errpl = 405, e.Error()
		case errors.Is(e, internal.ErrNotFound) || errors.Is(e, jujuerr.NotFound) || errors.Is(e, jujuerr.UserNotFound):
			errcode, errpl = 404, e.Error()
		case errors.Is(e, internal.ErrUnauthorized) || errors.Is(e, jujuerr.Unauthorized):
			errcode, errpl = 401, e.Error()
		case errors.Is(e, jujuerr.NotImplemented):
			errcode, errpl = 501, e.Error()
		}
	}

	return errcode, gin.H{`error`: errpl}
}
