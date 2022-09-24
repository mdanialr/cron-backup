package helper

import "fmt"

// LogStart append [START] and | prefix | to the given message and return it. Also append newline at the end.
func LogStart(prefix, message string) string {
	return fmt.Sprintf("[BEGN] | %s | %s\n", prefix, message)
}

// LogDone append [DONE] to the given message and return it.
func LogDone(prefix, message string) string {
	return fmt.Sprintf("[DONE] | %s | %s\n", prefix, message)
}
