package gobackup

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jeamon/gobackup/pkg/app"
)

// Execute is the entry point of the application. It processes the command-line arguments
// and calls the associated routine (version or help or monitor or logs filtering) if valid.
func Execute(buildTime, commit, tag string) int {
	buildTime = normalizeFlag(buildTime)
	commit, tag = normalizeFlag(commit), normalizeFlag(tag)

	var option Option
	monitorCommand, logsCommand := option.SetFlags()

	if ok := printVersionOrHelp(os.Stdout, os.Args, buildTime, commit, tag); ok {
		return 0
	}

	if ok := isMonitorOrLogsViewCommandArgs(os.Args); !ok {
		fmt.Printf("Invalid syntax. Run '%s help' for usage.\n", filepath.Base(os.Args[0]))
		return 1
	}

	command := strings.ToLower(os.Args[1])
	switch command {
	case "monitor":
		if err := monitorCommand.Parse(os.Args[2:]); err != nil {
			log.Printf("app monitoring mode: failed to parse arguments provided: %v", err)
			return 1
		}

		exitCode, err := app.Backup(runtime.NumCPU()*2-1, option.logFilePath, option.srcPath, option.dstPath, commit, tag)
		if err != nil {
			log.Printf("app monitoring mode: %v", err)
		}
		return exitCode

	case "logs":
		if err := logsCommand.Parse(os.Args[2:]); err != nil {
			log.Printf("app logs filtering mode: failed to parse arguments provided: %v", err)
			return 1
		}

		exitCode, err := app.ViewLogs(option.fileToFilter, option.date, option.regex)
		if err != nil {
			log.Printf("app logs filtering mode: logs filtering mode: %v", err)
		}
		return exitCode
	}
	return 0
}

// printVersionOrHelp checks if the command line arguments aims to show the application
// version or help details. If yes, then outputs the rquested information on `out`.
func printVersionOrHelp(out io.Writer, args []string, buildTime, commit, tag string) bool {
	if len(args) != 2 {
		return false
	}

	if isVersionCommand(args[1]) {
		printVersion(out, buildTime, commit, tag)
		return true
	}

	if isHelpCommand(args[1]) {
		printHelp(out)
		return true
	}

	return false
}

// isMonitorOrLogsViewCommandArgs checks if the commands line arguments satisfy
// the minimal requirements to run the app into monitoring or log-filtering mode.
// To run the app we expect at least 5 arguments. See commands examples below :
// appExec monitor [-file <logpath>] -source <src> -backup <dst>
// appExec logs [-file <logpath>] -date <date> -regex <regex>
func isMonitorOrLogsViewCommandArgs(args []string) bool {
	if len(args) < 6 {
		return false
	}

	cmd := args[1]
	if cmd != "monitor" && cmd != "logs" {
		return false
	}
	return true
}

// isVersionCommand checks if argument is any `version` keyword.
func isVersionCommand(arg string) bool {
	arg = strings.ToUpper(arg)
	if arg == "VERSION" || arg == "--VERSION" || arg == "-V" {
		return true
	}
	return false
}

// isHelpCommand checks if argument is any `help` keyword.
func isHelpCommand(arg string) bool {
	arg = strings.ToUpper(arg)
	if arg == "HELP" || arg == "--HELP" || arg == "-H" {
		return true
	}
	return false
}

// printVersion outputs the version information into `out`.
func printVersion(out io.Writer, buildTime, commit, tag string) {
	fmt.Fprintf(out, "Version: %s\nGo version: %s\nGit commit: %s\nOS/Arch: %s/%s\nBuilt: %s\n",
		tag, runtime.Version(), commit, runtime.GOOS, runtime.GOARCH, buildTime)
}

// printVersion dumps the `usage` value into `out`.
func printHelp(out io.Writer) {
	fmt.Fprintf(out, "%s\n", usage)
}

// normalizeFlag returns `(unknown)` for empty flag.
func normalizeFlag(flag string) string {
	if flag == "" {
		return "(unknown)"
	}
	return flag
}
