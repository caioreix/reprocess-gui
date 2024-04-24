package utils

import (
	"runtime"
	"strings"
)

// GetCaller retrieves the name of the calling function.
// It skips 'skip' number of stack frames.
func GetCaller(skip int) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	parts := strings.Split(frame.Function, "/")
	return parts[len(parts)-1]
}
