package gobackup

import "flag"

// Option represents user inputs.
type Option struct {
	srcPath      string
	dstPath      string
	date         string
	regex        string
	logFilePath  string
	fileToFilter string
}

// SetFlags configures flags for both monitoring and logs filtering commands.
func (o *Option) SetFlags() (*flag.FlagSet, *flag.FlagSet) {
	monitorCommand := flag.NewFlagSet("monitor", flag.ExitOnError)
	monitorCommand.StringVar(&o.logFilePath, "file", "file.log", "path to the file for logging.")
	monitorCommand.StringVar(&o.srcPath, "source", "", "path of the source folder to monitor its content.")
	monitorCommand.StringVar(&o.dstPath, "backup", "", "path of the backup folder for storing copied files.")

	logsCommand := flag.NewFlagSet("logs", flag.ExitOnError)
	logsCommand.StringVar(&o.fileToFilter, "file", "file.log", "path to the log file for filtering.")
	logsCommand.StringVar(&o.date, "date", "", "date of log entries to display.")
	logsCommand.StringVar(&o.regex, "regex", "", "regex to match against filename into logs.")
	return monitorCommand, logsCommand
}
