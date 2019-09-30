#! /bin/bash
echo "starting build..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/tool.linux
#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/tool.exe
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dist/tool.osx
#CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o dist/tool.arm
ls -lhrta dist/
echo "finish build. please see below dist directory."