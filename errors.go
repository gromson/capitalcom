package capitalcom

import (
	"fmt"

	werrors "github.com/gromson/capitalcom/pkg/errors"
)

type RequestPayloadEncodingError struct{ werrors.WrapperError }

func NewRequestPayloadEncodingError(err error) RequestPayloadEncodingError {
	return RequestPayloadEncodingError{werrors.Wrap(err, "failed to encode request payload")}
}

type RequestCreationError struct{ werrors.WrapperError }

func NewRequestCreationError(err error) RequestCreationError {
	return RequestCreationError{werrors.Wrap(err, "failed to create an HTTP request")}
}

type HTTPRequestError struct{ werrors.WrapperError }

func NewHTTPRequestError(err error) HTTPRequestError {
	return HTTPRequestError{werrors.Wrap(err, "HTTP request error")}
}

type ResponsePayloadDecodingError struct{ werrors.WrapperError }

func NewResponsePayloadDecodingError(err error) ResponsePayloadDecodingError {
	return ResponsePayloadDecodingError{werrors.Wrap(err, "failed to decode response payload")}
}

type APIError struct {
	statusCode int
	errorCode  string
}

func NewAPIError(statusCode int, errorCode string) APIError {
	return APIError{
		statusCode: statusCode,
		errorCode:  errorCode,
	}
}

func (e APIError) Error() string {
	return fmt.Sprintf("API returned an error, statusCode: %d, errorCode: %s", e.statusCode, e.errorCode)
}

func (e APIError) StatusCode() int {
	return e.statusCode
}

func (e APIError) ErrorCode() string {
	return e.errorCode
}
