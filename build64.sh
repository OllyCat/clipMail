#!/bin/sh

GOOS=windows CGO_ENABLED=1 CXX=x86_64-w64-mingw32-cpp-win32 CC=x86_64-w64-mingw32-gcc-win32 go build -v -o clipMail.x64.exe -ldflags -H=windowsgui
