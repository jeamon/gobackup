package main

import (
	"os"
	"runtime"

	"github.com/jeamon/gobackup/cmd/gobackup"
)

var (
	// holds the datetime when this executable was built.
	BuildTime string
	// holds latest git commit of codebase used to build the executable.
	GitCommit string
	// holds latest git tag of codebase used to build the executable.
	GitTag string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	code := gobackup.Execute(BuildTime, GitCommit, GitTag)
	os.Exit(code)
}
