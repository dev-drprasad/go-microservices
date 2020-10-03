package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc/metadata"
)

type key int

const (
	logKey key = iota
	authUserKey
)

// var Unauthorized = map[string]string{"message": "You are not authorized to access this resource"}
// var Internal = map[string]string{"message": "An error occured"}
// var BadRequest = map[string]string{"message": "Invalid payload"}
// var NotFound = map[string]string{"message": ""}

var ErrNoResourceFound = errors.New("no resource found")
var ErrFKViolation = errors.New("violates foreign key constraint")
var ErrViolation = errors.New("sql integrity violation")

func NewContextWithLogger(ectx echo.Context) context.Context {
	ctx := ectx.Request().Context()
	ctx = context.WithValue(ctx, "user", ectx.Get("user"))
	return context.WithValue(ctx, logKey, ectx.Logger())
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

type Float64 float64

func (v *Float64) UnmarshalJSON(data []byte) error {
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	switch temp.(type) {
	case float64:
		*v = Float64(temp.(float64))
	case int64:
		t := temp.(int64)
		*v = Float64(t)
	default:
		return fmt.Errorf("Expected int64 or float64 as cost. but got %T", temp)

	}

	return nil
}

type apiError struct {
	Message string `json:"message"`
}

func NewAPIError(message string) *apiError {
	return &apiError{Message: message}
}

type CountOnDate struct {
	Date  time.Time `json:"date"`
	Count uint      `json:"count"`
}

// https://stackoverflow.com/questions/12486436
func MakeInsertSQLWithValue(stmt string, value string, rows int, start int) string {
	// Combine the base SQL string and N value strings
	values := make([]string, 0, rows)
	for i := 0; i < rows; i++ {
		values = append(values, "("+value+")")
	}

	stmt = fmt.Sprintf(stmt, strings.Join(values, ","))

	// Convert all of the "?" to "$1", "$2", "$3", etc.
	// (which is the way that pgx expects query variables to be)
	numArgs := strings.Count(stmt, "?")
	stmt = strings.ReplaceAll(stmt, "?", "$%v")
	replaced := make([]interface{}, 0, rows)
	for i := start; i <= numArgs; i++ {
		replaced = append(replaced, strconv.Itoa(i))
	}
	return fmt.Sprintf(stmt, replaced...)
}

// MakeInsertSQL is a helper function to prepare a SQL query for a bulk insert
// MakeInsertSQL is used over getBulkInsertSQL when all of the values are plain question
// marks (e.g. a 1-for-1 value insertion)
// The example given for getBulkInsertSQL is such a query
func MakeInsertSQL(stmt string, perRow int, rows int) string {
	placeholders := make([]string, 0, perRow)
	for i := 0; i < perRow; i++ {
		placeholders = append(placeholders, "?")
	}
	value := strings.Join(placeholders, ", ")
	return MakeInsertSQLWithValue(stmt, value, rows, 1)
}

func CheckSQLViolation(err error) error {
	// https://www.postgresql.org/docs/13/errcodes-appendix.html
	codes := []string{"23000", "23001", "23502", "23503", "23505", "23514"}
	for _, code := range codes {
		if strings.Contains(err.Error(), code) {
			return ErrViolation
		}
	}
	return err
}
