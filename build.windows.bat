@echo off
git rev-parse --short HEAD > tempFile && set /P gitCommit=<tempFile
git describe --tags --abbrev=0 > tempFile && set /P gitTag=<tempFile
for /f "tokens=*" %%i in ('tzutil /g') do echo %date% %time:~0,-3% %%i > tempFile && set /P buildTime=<tempFile
del tempFile
go build -o bin/gobackup.exe -a -ldflags "-extldflags '-static' -X 'main.GitCommit=%gitCommit%' -X 'main.GitTag=%gitTag%' -X 'main.BuildTime=%buildTime%'" main.go