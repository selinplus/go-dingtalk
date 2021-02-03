.PHONY: all init lnx win mac check run clean help

BINARY=go-dingtalk

all: win run

init: win
	./${BINARY}.exe init

lnx:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build

run:
	./${BINARY}.exe

check:
	go fmt ./...
	go vet ./...

clean:
	go clean -i .

help:
	@echo "make all - doc, win, run"
	@echo "make lnx - compile Go code, generate Linux binary file"
	@echo "make mac - compile Go code, generate Mac binary file"
	@echo "make win - compile Go code, generate Windows executable file"
	@echo "make run - run executable file"
	@echo "make clean - remove object files and cached files"
	@echo "make check - run go tool 'fmt' and 'vet'"