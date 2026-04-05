package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Response is the standard JSON response envelope
type Response struct {
	OK        bool        `json:"ok"`
	Command   string      `json:"command"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorBody  `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// ErrorBody is the error payload
type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// StreamEvent is a single line in --watch streaming mode
type StreamEvent struct {
	Type           string      `json:"type"` // "status", "progress", "done", "error"
	Message        string      `json:"message,omitempty"`
	Current        int         `json:"current,omitempty"`
	Total          int         `json:"total,omitempty"`
	DurationSecs   float64     `json:"duration_seconds,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Timestamp      string      `json:"timestamp"`
}

// Success returns a successful JSON response
func Success(command string, data interface{}) *Response {
	return &Response{
		OK:        true,
		Command:   command,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// ErrorResponse returns an error JSON response
func ErrorResponse(command string, err error) *Response {
	code := "UNKNOWN"
	if e, ok := err.(Error); ok {
		code = e.Code
	}
	return &Response{
		OK:      false,
		Command: command,
		Error: &ErrorBody{
			Code:    code,
			Message: err.Error(),
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Print writes a response to stdout as JSON
func Print(command string, data interface{}) {
	resp := Success(command, data)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "")
	enc.Encode(resp)
}

// PrintError writes an error response to stderr
func PrintError(command string, err error) {
	resp := ErrorResponse(command, err)
	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "")
	enc.Encode(resp)
}

// Error is a typed error with a code
type Error struct {
	Code    string
	Message string
	Details interface{}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewError creates a typed error
func NewError(code, message string, details interface{}) error {
	return Error{Code: code, Message: message, Details: details}
}

// Stream prints a streaming event line to stdout
func Stream(t, message string, args ...interface{}) {
	ev := StreamEvent{
		Type:      t,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	if len(args) > 0 {
		ev.Message = fmt.Sprintf(message, args...)
	}
	json.NewEncoder(os.Stdout).Encode(ev)
}

// StreamWithData prints a streaming event with data
func StreamWithData(t string, data interface{}) {
	ev := StreamEvent{
		Type:      t,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	json.NewEncoder(os.Stdout).Encode(ev)
}
