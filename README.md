# gobackup

This is a cross-platform go project to demonstrate the usage of my [gorsn](github.com/jeamon/gorsn) package and how useful it can help
to live monitor a directory and its sub-directories and react to any file changes.


## Setup

On Windows, Linux macOS, and FreeBSD you will be able to download the pre-built binaries once available.
If your system has [Go >= 1.21.1](https://golang.org/dl/) you can pull the codebase and build from the source.

```
# on windows
git clone https://github.com/jeamon/gobackup.git && cd gobackup
make build (or use `go build -o gobackup.exe main.go`)

# on linux and others
git clone https://github.com/jeamon/gobackup.git && cd gobackup
make build
```


## Usage


```Usage:
    
	This tool allows to monitor a hot source folder and backup any regular file created or modified
	inside this folder and its sub-folders. Use Ctrl-C to stop the program. Before it exits, the
	backup folder content will be saved into a zip archive using the datetime and process id into
	the filename. Finally it allows you to view logs entries based on the date and filename regex.
	Specify the path towards the log file for filtering. If not specified it default to <file.log>.
	Use CTRL+C to stop the program on windows machines. On Linux and MacOS you can use Kill command. 
	
	gobackup [version | help ]
	gobackup monitor -source <path-to-hot-folder> -backup <path-to-backup-folder>
	gobackup logs -file <logfile-path> -date <yyyy-mm-dd> -regex <filename-regex>

    Examples:
	
	$ ./gobackup monitor -source "C:\demo\source" -backup "C:\demo\backup"
	$ ./gobackup logs -date 2023-08-14 -regex *.bak
	$ ./gobackup logs -file file.log -date 2023-08-14 -regex *.bak
	
	$ ./gobackup help
	$ ./gobackup version
	
```


## License

please check & read [the license details](https://github.com/jeamon/gobackup/blob/main/LICENSE) or [reach out to me](https://blog.cloudmentor-scale.com/contact) before any action.