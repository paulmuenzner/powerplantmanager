package responsehandler

import (
	"encoding/json"
	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"
	"net/http"
)

// CustomError represents a custom error with a message and HTTP status code.
type CustomError struct {
	Message    string
	StatusCode int
}

// HandleSuccess handles a success and writes the response with the appropriate status code.
// 20 HTTP status codes
func HandleSuccess(w http.ResponseWriter, message string, status SuccessStatus, data ...interface{}) {
	var statusCode int
	var statusName string

	// Recover from panics in this function
	defer func() {
		if r := recover(); r != nil {
			// Log or handle the panic here
			logger.GetLogger().Error("Recovered from panic in HandleSuccess:", r)
			// Optionally, you can send an HTTP 500 Internal Server Error response
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	switch status {
	case OK:
		statusCode = http.StatusOK // 200
		statusName = "OK"
	case Created:
		statusCode = http.StatusCreated // 201
		statusName = "Created"
	case Accepted:
		statusCode = http.StatusAccepted // 202
		statusName = "Accepted"
	case NonAuthoritativeInfo:
		statusCode = http.StatusNonAuthoritativeInfo // 203
		statusName = "NonAuthoritativeInfo"
	case NoContent:
		statusCode = http.StatusNoContent // 204
		statusName = "NoContent"
	case ResetContent:
		statusCode = http.StatusResetContent // 205
		statusName = "ResetContent"
	case PartialContent:
		statusCode = http.StatusPartialContent // 206
		statusName = "PartialContent"
	case MultiStatus:
		statusCode = http.StatusMultiStatus // 207
		statusName = "MultiStatus"
	case AlreadyReported:
		statusCode = http.StatusAlreadyReported // 226
		statusName = "AlreadyReported"
	case IMUsed:
		statusCode = http.StatusIMUsed // 229
		statusName = "IMUsed"
	default:
		statusCode = http.StatusOK // 200
		statusName = "OK"
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Set the response status code
	w.WriteHeader(statusCode)

	// Content response
	successResponse := map[string]interface{}{
		"result": "success",
		"status": map[string]interface{}{
			"statusCode": statusCode,
			"statusName": statusName,
		},
		"message": message,
		"data":    data,
	}

	// Encode the success response map to JSON
	encodedResponse, err := json.Marshal(successResponse)
	if err != nil {
		// Handle the error, e.g., log it or return an HTTP 500 Internal Server Error
		logger.GetLogger().Error("Error in responseHandler. Cannot encode the success response map to JSON. Error: ", err)
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	// Write the encoded JSON to the response writer
	_, err = w.Write(encodedResponse)
	if err != nil {
		// Handle the error, e.g., log it or return an HTTP 500 Internal Server Error
		logger.GetLogger().Error("Error in responseHandler. Cannot write the encoded JSON to the response writer. Error: ", err)
		http.Error(w, "Error writing JSON response", http.StatusInternalServerError)
		return
	}

	// Check if the response writer supports flushing
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

}
