.PHONY: build-mac, build-linux, build-windows-64, build-windows-32, build-all

build-linux:
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./out/healthcheckLINUX main.go
	@echo "[OK] Files build to linux"

build-windows-64:
	@GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./out/healthcheck64.exe main.go
	@echo "[OK] Files build to windows(64bit)"

build-windows-32:
	@GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o ./out/healthcheck32.exe main.go
	@echo "[OK] Files build to windows(32bit)"

build-mac:
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ./out/healthcheckMAC main.go
	@echo "[OK] Files build to OSX "

build-all: build-linux build-windows-64 build-windows-32 build-mac
