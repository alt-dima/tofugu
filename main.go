package main

import "github.com/alt-dima/tofugu/cmd"

var (
	version string = "undefined"
)

func main() {
	cmd.SetVersionInfo(version)
	cmd.Execute()
}
