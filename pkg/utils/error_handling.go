package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/getsentry/sentry-go"
)

const (
	// HTTP status codes were copied from https://en.wikipedia.org/wiki/List_of_HTTP_status_codes
	StatusInvalidToken = 498
	// Custom status code errors
	ErrCodeQueryError = 1001
)

var (
	ErrQueryFailed = errors.New("query failed")
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func HandleErrors(ctx context.Context, errs ...error) {
	if len(errs) > 0 {
		for _, err := range errs {
			_, fn, line, _ := runtime.Caller(1)
			fns := strings.Split(fn, "/")
			log.Println(
				"[ERROR]",
				"txn_id:",
				ctx.Value("requestid"),
				"|",
				"file:",
				fns[len(fns)-1],
				"|",
				fmt.Sprintf("line: %d", line),
				"|",
				err.Error(),
			)

			sentry.CaptureException(err)
		}
	}
}
