package main

import (
	"os"

	"github.com/fnando/bolt/common"
	bolt "github.com/fnando/bolt/internal"
)

func main() {
	workingDir, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()

	exitcode := bolt.Run(
		workingDir,
		homeDir,
		common.Output{Stdout: os.Stdout, Stderr: os.Stderr},
	)

	os.Exit(exitcode)
}
