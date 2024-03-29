package router

import (
	"bytes"
	"fmt"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"runtime/debug"
	"sdkeji/visitor/pkg/code"
	"sort"
	"strings"
)

func routingLogger(logger zerolog.Logger) routing.Handler {
	return access.CustomLogger(func(req *http.Request, rw *access.LogResponseWriter, elapsed float64) {
		clientIP := access.GetClientIP(req)
		logger.Info().
			Str("ip", clientIP).
			Str("proto", req.Proto).
			Str("method", req.Method).
			Str("url", req.URL.String()).
			Int("status", rw.Status).
			Int64("size", rw.BytesWritten).
			Float64("duration", elapsed).
			Msg("access log.")
	})
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func errorHandler(logger zerolog.Logger) routing.Handler {
	return func(c *routing.Context) (err error) {
		defer func() {
			if e := recover(); e != nil {
				logger.Error().Msgf("recovered from panic: %v", e)
				fmt.Print(string(debug.Stack()))
				sendError(logger, c, code.NewAPIError(code.InternalServerError), http.StatusInternalServerError)
				c.Abort()
				err = nil
			}
		}()
		err = c.Next()
		if err != nil {
			c.Abort()

			if err, ok := err.(stackTracer); ok {
				buf := new(bytes.Buffer)
				buf.WriteString(fmt.Sprintf("error with stacktrace returned: %v\n", err))
				for _, f := range err.StackTrace() {
					buf.WriteString(fmt.Sprintf("%+v\n", f))
				}
				fmt.Fprint(os.Stderr, buf.String())
			}

			switch e := errors.Cause(err).(type) {
			case validation.Errors:
				type validationError struct {
					Field string `json:"field"`
					Error string `json:"error"`
				}
				result := make([]validationError, 0)
				fields := make([]string, 0)
				for field := range e {
					fields = append(fields, field)
				}
				sort.Strings(fields)
				for _, field := range fields {
					err := e[field]
					result = append(result, validationError{
						Field: field,
						Error: err.Error(),
					})
				}
				apiErr := code.NewAPIError(code.InvalidData).WithDetails(result)
				sendError(logger, c, apiErr, apiErr.StatusCode())
			case code.APIError:
				sendError(logger, c, e, e.StatusCode())
			case routing.HTTPError:
				sendError(logger, c, code.NewAPIError(code.InvalidData).WithMessage(e.Error()), e.StatusCode())
			default:
				logger.Error().Err(err).Msg("unknown error.")
				sendError(logger, c, code.NewAPIError(code.InternalServerError), http.StatusInternalServerError)
				return nil
			}
		}
		return nil
	}
}

func notFound(c *routing.Context) error {
	methods := c.Router().FindAllowedMethods(c.Request.URL.Path)
	if len(methods) == 0 {
		return code.NewAPIError(code.NotFound)
	}
	methods["OPTIONS"] = true
	ms := make([]string, len(methods))
	i := 0
	for method := range methods {
		ms[i] = method
		i++
	}
	sort.Strings(ms)
	c.Response.Header().Set("Allow", strings.Join(ms, ", "))
	if c.Request.Method != "OPTIONS" {
		return code.NewAPIError(code.MethodNotAllowed)
	}
	c.Abort()
	return nil
}

func sendError(logger zerolog.Logger, c *routing.Context, err error, status int) {
	c.Response.WriteHeader(status)
	c.Response.Header().Set("X-Content-Type-Options", "nosniff")
	err = c.Write(err)
	if err != nil {
		logger.Error().Err(err).Msg("fail to write error.")
	}
}
