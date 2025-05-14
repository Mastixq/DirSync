package logger

import (
	"fmt"
	"os"
)

func Info(msg string, args ...any) {
	fmt.Printf("[INFO] "+msg+"\n", args...)
}

func Error(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "[ERROR] "+msg+"\n", args...)
}
