package logger

import "log"

func Info(message string, args ...any) {
	log.Println(append([]any{message}, args...)...)
}
