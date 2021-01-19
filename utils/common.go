package utils

import (
	"fmt"
	"os"
)

func QuitMsg(msg string) {
	fmt.Fprintf(os.Stderr, "%s", msg)
	os.Exit(1)
}
