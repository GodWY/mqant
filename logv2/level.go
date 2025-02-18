package logv2

import "strings"

// Level is a log level.
type Level int8

// LevelKey is log level key.
const LevelKey = "level"

const (
	// LevelDebug is log debug level.
	LevelDebug Level = iota - 1
	// LevelInfo is log info level.
	LevelInfo
	// LevelWarn is log warn level.
	LevelWarn
	// LevelError is log error level.
	LevelError
	// LevelFatal is log fatal level
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// ParseLevel parses a level string into a log Level value.
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	}
	return LevelInfo
}
