package log

import (
	"context"

	"github.com/labstack/echo/v4"
	log "github.com/labstack/gommon/log"
)

type logKey string

const lk logKey = "log"

func NewContext(ectx echo.Context) context.Context {
	ctx := ectx.Request().Context()
	ctx = context.WithValue(ctx, "user", ectx.Get("user"))
	return context.WithValue(ctx, lk, ectx.Logger())
}

func FromContext(ctx context.Context) echo.Logger {
	var logger echo.Logger
	var ok bool
	if logger, ok = ctx.Value(lk).(echo.Logger); !ok {
		logger = log.New("")
	}
	return logger
}

func LoggerWithRequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rID := c.Response().Header().Get(echo.HeaderXRequestID)
		c.Logger().SetLevel(log.DEBUG)
		c.Logger().SetPrefix(rID)

		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}
