.PHONY:
.SILENT:

build:
	go mod download && go build -o ./.bin/app ./cmd/app/main.go
