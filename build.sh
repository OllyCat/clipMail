#!/bin/sh

GOOS=windows GOARCH=386 CGO_ENABLED=1 CXX=i686-w64-mingw32-cpp-win32 CC=i686-w64-mingw32-gcc-win32 go build -v -o clipMail.386.exe -ldflags -H=windowsgui
