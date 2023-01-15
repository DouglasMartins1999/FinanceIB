#!/bin/sh
GOOS=linux GOARCH=386 go build -o dist/financeIB_linux32
GOOS=linux GOARCH=amd64 go build -o dist/financeIB_linux64
GOOS=linux GOARCH=arm64 go build -o dist/financeIB_linuxARM
GOOS=windows GOARCH=386 go build -o dist/financeIB_win32.exe -ldflags -H=windowsgui
GOOS=windows GOARCH=amd64 go build -o dist/financeIB_win64.exe -ldflags -H=windowsgui
GOOS=darwin GOARCH=arm64 go build -o dist/financeIB_macOS_silicon.app
GOOS=darwin GOARCH=amd64 go build -o dist/financeIB_macOS_intel.app