package main

import (
	"os"
	"path/filepath"

	"github.com/ashellwig/ashhellwig-go/pkg/cmd"
)

func main() {
	baseName := filepath.Base(os.Args[0]).Execute()
	cmd.CheckError()
}
