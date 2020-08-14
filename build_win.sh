#!/bin/bash
sudo apt update
sudo apt install gcc-mingw-w64-x86-64

CGO_ENABLED=1 
GOOS=windows 
GOARCH=amd64 
CC=x86_64-w64-mingw32-gcc
go build -o builds/$1/casgo_win64.exe app.go
