#!/bin/bash
sudo apt update
sudo apt install gcc-arm-linux-gnueabi gcc-arm-linux-gnueabihf gcc-aarch64-linux-gnu gcc-mipsel-linux-gnu

#windows x64
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o builds/$1/casgo_win64.exe app.go
#linux amd64
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o builds/$1/casgo_linux_amd64 app.go
#linux_arm_eabi
CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabi-gcc go build -o builds/$1/casgo_linux_armeabi app.go
#linux_arm_eabihf
CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc go build -o builds/$1/casgo_linux_armeabihf app.go
#linux_aarch64
CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=aarch64-linux-gnu-gcc go build -o builds/$1/casgo_linux_aarch64 app.go
#linux_mipsle
CGO_ENABLED=1 GOOS=linux GOARCH=mips32le CC=mipsel-linux-gnu-gcc go build -o builds/$1/casgo_linux_mipsel app.go



