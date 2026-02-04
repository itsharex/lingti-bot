package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// Level represents the logging level
type Level int

const (
	// LevelSilent shows minimal output
	LevelSilent Level = iota
	// LevelInfo shows command execution (default)
	LevelInfo
	// LevelVerbose shows command results
	LevelVerbose
	// LevelVeryVerbose shows debug messages
	LevelVeryVerbose
)

var (
	currentLevel = LevelInfo
	mu           sync.RWMutex

	// Loggers for different levels
	silentLogger      *log.Logger
	infoLogger        *log.Logger
	verboseLogger     *log.Logger
	veryVerboseLogger *log.Logger
)

func init() {
	silentLogger = log.New(os.Stderr, "", 0)
	infoLogger = log.New(os.Stderr, "", log.LstdFlags)
	verboseLogger = log.New(os.Stderr, "[VERBOSE] ", log.LstdFlags)
	veryVerboseLogger = log.New(os.Stderr, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
}

// ParseLevel parses a string into a Level
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(s) {
	case "silent", "s", "0":
		return LevelSilent, nil
	case "info", "i", "1":
		return LevelInfo, nil
	case "verbose", "v", "2":
		return LevelVerbose, nil
	case "very-verbose", "vv", "debug", "d", "3":
		return LevelVeryVerbose, nil
	default:
		return LevelInfo, fmt.Errorf("unknown log level: %s (use: silent, info, verbose, very-verbose)", s)
	}
}

// SetLevel sets the global log level
func SetLevel(level Level) {
	mu.Lock()
	defer mu.Unlock()
	currentLevel = level

	// Disable standard log output for silent mode
	if level == LevelSilent {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}
}

// GetLevel returns the current log level
func GetLevel() Level {
	mu.RLock()
	defer mu.RUnlock()
	return currentLevel
}

// Silent logs a message at silent level (always shown except errors)
func Silent(format string, v ...any) {
	silentLogger.Printf(format, v...)
}

// Info logs a message at info level
func Info(format string, v ...any) {
	mu.RLock()
	level := currentLevel
	mu.RUnlock()

	if level >= LevelInfo {
		infoLogger.Printf(format, v...)
	}
}

// Verbose logs a message at verbose level
func Verbose(format string, v ...any) {
	mu.RLock()
	level := currentLevel
	mu.RUnlock()

	if level >= LevelVerbose {
		verboseLogger.Printf(format, v...)
	}
}

// Debug logs a message at very-verbose/debug level
func Debug(format string, v ...any) {
	mu.RLock()
	level := currentLevel
	mu.RUnlock()

	if level >= LevelVeryVerbose {
		veryVerboseLogger.Printf(format, v...)
	}
}

// Error always logs errors regardless of level
func Error(format string, v ...any) {
	silentLogger.Printf("[ERROR] "+format, v...)
}

// IsVerbose returns true if verbose logging is enabled
func IsVerbose() bool {
	mu.RLock()
	defer mu.RUnlock()
	return currentLevel >= LevelVerbose
}

// IsDebug returns true if debug logging is enabled
func IsDebug() bool {
	mu.RLock()
	defer mu.RUnlock()
	return currentLevel >= LevelVeryVerbose
}

// IsSilent returns true if silent mode is enabled
func IsSilent() bool {
	mu.RLock()
	defer mu.RUnlock()
	return currentLevel == LevelSilent
}
