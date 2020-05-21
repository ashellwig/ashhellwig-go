package cmd

import (
	"context"
	"fmt"
	"os"
)

// CheckError ensures there are no runtime errors when using the application.
func CheckError(err error) {
	if err != nil {
		if err != context.Canceled {
			fmt.Fprintf(os.Stderr, "An error ocurred: %v\n", err)
		}
		os.Exit(1)
	}
}

// Exit prints an error message and terminates the program.
func Exit(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
