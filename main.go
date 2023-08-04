package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nyaosorg/nyagos/internal/defined"
	"github.com/nyaosorg/nyagos/internal/frame"
	"github.com/nyaosorg/nyagos/internal/mains"
	"github.com/nyaosorg/nyagos/internal/onexit"
)

var version string

func main() {
	defer onexit.Done()

	frame.Version = strings.TrimSpace(version)
	if err := frame.Start(mains.Main); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err)
		defer os.Exit(1)
	}
	if defined.DBG {
		os.Stdin.Read(make([]byte, 1))
	}
}
