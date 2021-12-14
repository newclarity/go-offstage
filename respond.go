package offstage

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mikeschinkel/go-only"
	"github.com/pkg/errors"
	"net/http"
)

// RespondOnSuccess responds via HTTP(S) with a JSON body in response
// on a successful API call. It return an error object if the JSON
// encoding fails for some reason.
func RespondOnSuccess(ctx echo.Context, data interface{}) error {
	return RespondOnSuccessWithStatusAndHeaders(ctx, http.StatusOK, data, nil)
}

// RespondOnSuccessWithHeaders is RespondOnSuccess but adds a parameter for HTTP headers
//goland:noinspection GoUnusedExportedFunction
func RespondOnSuccessWithHeaders(ctx echo.Context, data interface{}, headers http.Header) error {
	return RespondOnSuccessWithStatusAndHeaders(ctx, http.StatusOK, data, headers)
}

// RespondOnSuccessWithStatus is RespondOnSuccessWith but adds a parameter for statuc
//goland:noinspection GoUnusedExportedFunction
func RespondOnSuccessWithStatus(ctx echo.Context, status int, data interface{}) error {
	return RespondOnSuccessWithStatusAndHeaders(ctx, status, data, nil)
}

// RespondOnSuccessWithStatusAndHeaders responds via HTTP(S) with a JSON body in response
// on a successful API call. It return an error object if the JSON
// encoding fails for some reason.
func RespondOnSuccessWithStatusAndHeaders(ctx echo.Context, status int, data interface{}, headers http.Header) error {
	setHeaders(ctx, headers)
	return ctx.JSON(status, data)
}

type ErrorResponse interface {
	GetMessage() string
	SetMessage(string)
	GetCode() string
	GetStatus() int32
	GetData() string
}

// RespondOnFailure responds via HTTP(S) with a JSON body in response
// on a failed API call. The JSON will contain an HTTP status code
// and a message describing the issue from the passed error instance.
func RespondOnFailure(ctx echo.Context, er ErrorResponse) (err error) {
	return RespondOnFailureWithHeaders(ctx, nil, er)
}

// RespondOnFailureWithHeaders responds via HTTP(S) with a JSON body in response
// on a failed API call. The JSON will contain an HTTP status code
// and a message describing the issue from the passed error instance.
func RespondOnFailureWithHeaders(ctx echo.Context, headers http.Header, er ErrorResponse) (err error) {
	for range only.Once {
		if er.GetMessage() == "" {
			err = wrapError("Message property cannot be empty", err)
		}
		if er.GetCode() == "" {
			err = wrapError("Code cannot be empty", err)
		}
		if er.GetStatus() == 0 {
			err = wrapError("Status cannot be empty", err)
			er.SetMessage(fmt.Sprintf("Status cannot be zero when calling RespondOnFailure(); %s",
				er.GetMessage()))
		}
		if err != nil {
			err = wrapError("error when calling RespondOnFailure()", err)
			break
		}
		if headers != nil {
			setHeaders(ctx, headers)
		}
		err = ctx.JSON(int(er.GetStatus()), er)
	}
	return err
}

func wrapError(msg string, err error) error {
	for range only.Once {
		if err == nil {
			err = errors.New(msg)
			break
		}
		err = fmt.Errorf(msg+"; %w", err)
	}
	return err
}

func setHeaders(ctx echo.Context, headers http.Header) {
	for range only.Once {
		if headers == nil {
			break
		}
		for k, _v := range headers {
			for _, v := range _v {
				ctx.Response().Header().Set(k, v)
			}
		}
	}
}
