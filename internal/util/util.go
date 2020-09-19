package util

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc/metadata"
)

type key int

const (
	logKey key = iota
)

func NewContextWithLogger(ctx echo.Context) context.Context {
	return context.WithValue(context.Background(), logKey, ctx.Logger())
}

func GetLoggerFromContext(ctx context.Context) echo.Logger {
	l, _ := ctx.Value(logKey).(echo.Logger)
	return l
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

func NewContextFromMetadata(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	rID := md["rid"][0]
	logger := log.New(rID)
	logger.SetLevel(log.DEBUG)
	return context.WithValue(ctx, logKey, logger)
}
