package jshapi

import (
	"net/http"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/go-stdlogger"
	"golang.org/x/net/context"
)

/*
Sender is a function type definition that allows consumers to customize how they
send and log API responses.
*/
type Sender func(context.Context, http.ResponseWriter, *http.Request, jsh.Sendable)

/*
DefaultSender is the default sender that will log 5XX errors that it encounters
in the process of sending a response.
*/
func DefaultSender(logger std.Logger) Sender {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, sendable jsh.Sendable) {
		sendableError, isType := sendable.(jsh.ErrorType)
		if isType && sendableError.StatusCode() >= 500 {
			logger.Printf("Returning ISE: %s\n", sendableError.Error())
		}

		sendError := jsh.Send(w, r, sendable)
		if sendError != nil && sendError.Status >= 500 {
			logger.Printf("Error sending response: %s\n", sendError.Error())
		}
	}
}
