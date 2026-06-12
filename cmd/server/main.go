// Command server is an alternate entry point for TinyBrain kept for
// compatibility with the documented install path:
//
//	go install github.com/rainmana/tinybrain/cmd/server@latest
//
// It is identical to cmd/tinybrain; the binary is just named "server".
package main

import "github.com/rainmana/tinybrain/v3/internal/app"

func main() {
	app.Main()
}
