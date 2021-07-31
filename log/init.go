package log

import (
	"fmt"
	"strings"
	"runtime"
)

func getFrame(skip int) runtime.Frame {
	targetIndex := skip + 2

	counters := make([]uintptr, targetIndex+2)
	n := runtime.Callers(0, counters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(counters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetIndex; frameIndex++ {
			var candidate runtime.Frame
			candidate, more = frames.Next()
			if frameIndex == targetIndex {
				frame = candidate
			}
		}
	}

	return frame
}

func log(module string, level string, message interface{}) {
	// Cleanup the module name
	parts := strings.Split(module, "/")
	name := parts[len(parts) -1]
	fmt.Printf("[%s]: %s: %s\n", level, name, fmt.Sprint(message))
}

func Info(message interface{}) {
	frame := getFrame(1)
	log(frame.Function, "INFO", message)
}

func Warn(message interface{}) {
	frame := getFrame(1)
	log(frame.Function, "WARN", message)
}

func Error(message interface{}) {
	frame := getFrame(1)
	log(frame.Function, "ERR", message)
}
