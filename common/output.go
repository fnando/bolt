package common

import "io"

type Output struct {
	Stdout io.Writer
	Stderr io.Writer
}
