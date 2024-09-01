package gobackup

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsVersionCommand(t *testing.T) {
	cases := []struct {
		arg  string
		want bool
	}{
		{"help", false},
		{"version", true},
		{"VERSION", true},
		{"--version", true},
		{"--VERSION", true},
		{"-v", true},
		{"-V", true},
		{"--v", false},
		{"-version", false},
	}

	for _, tc := range cases {
		t.Run(tc.arg, func(t *testing.T) {
			got := isVersionCommand(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsHelpCommand(t *testing.T) {
	cases := []struct {
		arg  string
		want bool
	}{
		{"version", false},
		{"help", true},
		{"HELP", true},
		{"--help", true},
		{"--HELP", true},
		{"-h", true},
		{"-H", true},
		{"--h", false},
		{"-help", false},
	}

	for _, tc := range cases {
		t.Run(tc.arg, func(t *testing.T) {
			got := isHelpCommand(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPrintVersionOrHelp(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want bool
	}{
		{
			"version and help", []string{"version", "help"}, false,
		},
		{
			"version", []string{"version"}, true,
		},
		{
			"--version", []string{"--version"}, true,
		},
		{
			"-v", []string{"-v"}, true,
		},
		{
			"help", []string{"help"}, true,
		},
		{
			"--help", []string{"--help"}, true,
		},
		{
			"-h", []string{"-h"}, true,
		},
		{
			"monitor", []string{"monitor"}, false,
		},
		{
			"logs", []string{"logs"}, false,
		},
		{
			"monitor command", strings.Fields("monitor -source srcpath -backup dstpath"), false,
		},
		{
			"logs viewer command", strings.Fields("logs -date date -regex *.zip"), false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"gobackup"}, tc.args...)
			got := printVersionOrHelp(io.Discard, args, "", "", "")
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsMonitorOrLogsViewCommandArgs(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want bool
	}{
		{
			"version",
			[]string{"version"},
			false,
		},
		{
			"help",
			[]string{"help"},
			false,
		},
		{
			"monitor",
			[]string{"monitor"},
			false,
		},
		{
			"logs",
			[]string{"logs"},
			false,
		},
		{
			"monitor shortest command",
			strings.Fields("monitor -source srcpath -backup dstpath"),
			true,
		},
		{
			"monitor longest command",
			strings.Fields("monitor -file log.txt -source srcpath -backup dstpath"),
			true,
		},
		{
			"logs viewer shortest command",
			strings.Fields("logs -date date -regex *.zip"),
			true,
		},
		{
			"logs viewer longest command",
			strings.Fields("logs -file log.txt -date date -regex *.bak"),
			true,
		},
		{
			"invalid monitor command",
			strings.Fields("monitor -source -backup dstpath"),
			false,
		},
		{
			"invalid logs viewer command",
			strings.Fields("logs -date -regex delete_"),
			false,
		},
		{
			"expected number of monitor command",
			strings.Fields("monitor -file file.log -source srcpath"),
			true,
		},
		{
			"expected number of log viewer command",
			strings.Fields("logs -file file.log -date date"),
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"gobackup"}, tc.args...)
			got := isMonitorOrLogsViewCommandArgs(args)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPrintHelp(t *testing.T) {
	out := bytes.NewBuffer(nil)
	printHelp(out)
	assert.Equal(t, usage, strings.TrimSuffix(out.String(), "\n"))
}

func TestNormalizeFlag(t *testing.T) {
	unknown := "(unknown)"
	got := normalizeFlag("")
	assert.Equal(t, unknown, got)
	got = normalizeFlag("flag")
	assert.Equal(t, "flag", got)
}

func TestExecute(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	cases := []struct {
		name string
		args []string
		want int
	}{
		{
			"version",
			[]string{"version"},
			0,
		},
		{
			"help",
			[]string{"help"},
			0,
		},
		{
			"monitor: missing arguments",
			[]string{"monitor"},
			1,
		},
		{
			"logs: missing arguments",
			[]string{"logs"},
			1,
		},
		{
			"monitor: shortest command with invalid paths",
			strings.Fields("monitor -source srcpath -backup dstpath"),
			1,
		},
		{
			"monitor: longest command with invalid paths",
			strings.Fields("monitor -file log.txt -source srcpath -backup dstpath"),
			1,
		},
		{
			"logs: shortest command with invalid date",
			strings.Fields("logs -date date -regex *.zip"),
			1,
		},
		{
			"logs: longest command with inexistant log file",
			strings.Fields("logs -file noexist.log.txt -date date -regex *.bak"),
			1,
		},
		{
			"unknown command",
			strings.Fields("unknown.command -date date -regex *.bak"),
			1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = append([]string{os.Args[0]}, tc.args...)
			got := Execute("", "", "")
			assert.Equal(t, tc.want, got)
		})
	}
}
