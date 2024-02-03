package errorhandler

import (
	"encoding/json"
	"net/http"
)

// CustomError represents a custom error with a message and HTTP status code.
type CustomError struct {
	Message    string
	StatusCode int
}

// HandleError handles a custom error and writes the response with the appropriate status code.
// 32 HTTP status codes
func HandleError(w http.ResponseWriter, message string, status CustomStatus) {
	var statusCode int
	var statusName string

	switch status {
	case BadRequest:
		statusCode = http.StatusBadRequest // 400
		statusName = "BadRequest"
	case Unauthorized:
		statusCode = http.StatusUnauthorized // 401
		statusName = "Unauthorized"
	case PaymentRequired:
		statusCode = http.StatusPaymentRequired // 402
		statusName = "PaymentRequired"
	case Forbidden:
		statusCode = http.StatusForbidden // 403
		statusName = "Forbidden"
	case NotFound:
		statusCode = http.StatusNotFound // 404
		statusName = "NotFound"
	case MethodNotAllowed:
		statusCode = http.StatusMethodNotAllowed // 405
		statusName = "MethodNotAllowed"
	case NotAcceptable:
		statusCode = http.StatusNotAcceptable // 406
		statusName = "NotAcceptable"
	case ProxyAuthRequired:
		statusCode = http.StatusProxyAuthRequired // 407
		statusName = "ProxyAuthRequired"
	case RequestTimeout:
		statusCode = http.StatusRequestTimeout // 408
		statusName = "RequestTimeout"
	case Conflict:
		statusCode = http.StatusConflict // 409
		statusName = "Conflict"
	case Gone:
		statusCode = http.StatusGone // 410
		statusName = "Gone"
	case LengthRequired:
		statusCode = http.StatusLengthRequired // 411
		statusName = "LengthRequired"
	case PreconditionFailed:
		statusCode = http.StatusPreconditionFailed // 412
		statusName = "PreconditionFailed"
	case RequestEntityTooLarge:
		statusCode = http.StatusRequestEntityTooLarge // 413
		statusName = "RequestEntityTooLarge"
	case RequestURITooLong:
		statusCode = http.StatusRequestURITooLong // 414
		statusName = "RequestURITooLong"
	case UnsupportedMediaType:
		statusCode = http.StatusUnsupportedMediaType // 415
		statusName = "UnsupportedMediaType"
	case RequestedRangeNotSatisfiable:
		statusCode = http.StatusRequestedRangeNotSatisfiable // 416
		statusName = "RequestedRangeNotSatisfiable"
	case ExpectationFailed:
		statusCode = http.StatusExpectationFailed // 417
		statusName = "ExpectationFailed"
	case Teapot:
		statusCode = http.StatusTeapot // 418
		statusName = "Teapot"
	case MisdirectedRequest:
		statusCode = http.StatusMisdirectedRequest // 421
		statusName = "MisdirectedRequest"
	case UnprocessableEntity:
		statusCode = http.StatusUnprocessableEntity // 422
		statusName = "UnprocessableEntity"
	case Locked:
		statusCode = http.StatusLocked // 423
		statusName = "Locked"
	case FailedDependency:
		statusCode = http.StatusFailedDependency // 424
		statusName = "FailedDependency"
	case UpgradeRequired:
		statusCode = http.StatusUpgradeRequired // 426
		statusName = "UpgradeRequired"
	case PreconditionRequired:
		statusCode = http.StatusPreconditionRequired // 428
		statusName = "PreconditionRequired"
	case TooManyRequests:
		statusCode = http.StatusTooManyRequests // 429
		statusName = "TooManyRequests"
	case RequestHeaderFieldsTooLarge:
		statusCode = http.StatusRequestHeaderFieldsTooLarge // 431
		statusName = "RequestHeaderFieldsTooLarge"
	case UnavailableForLegalReasons:
		statusCode = http.StatusUnavailableForLegalReasons // 451
		statusName = "UnavailableForLegalReasons"
	case InternalServerError:
		statusCode = http.StatusInternalServerError // 500
		statusName = "InternalServerError"
	case NotImplemented:
		statusCode = http.StatusNotImplemented // 501
		statusName = "NotImplemented"
	case BadGateway:
		statusCode = http.StatusBadGateway // 502
		statusName = "BadGateway"
	case ServiceUnavailable:
		statusCode = http.StatusServiceUnavailable // 503
		statusName = "ServiceUnavailable"
	case GatewayTimeout:
		statusCode = http.StatusGatewayTimeout // 504
		statusName = "GatewayTimeout"
	case HTTPVersionNotSupported:
		statusCode = http.StatusHTTPVersionNotSupported // 505
		statusName = "HTTPVersionNotSupported"
	default:
		statusCode = http.StatusInternalServerError // 500
		statusName = "InternalServerError"
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Set the response status code
	w.WriteHeader(statusCode)

	// Content response
	errorResponse := map[string]interface{}{
		"result": "error",
		"status": map[string]interface{}{
			"statusCode": statusCode,
			"statusName": statusName,
		},
		"message": message,
	}

	// Encode the error response map to JSON and write it to the response writer
	json.NewEncoder(w).Encode(errorResponse)

}
