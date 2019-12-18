RELEASE=-s -w
UPXBIN=/usr/local/bin/upx
GOBIN=/usr/local/bin/go
GOOS=$(shell uname -s | tr [A-Z] [a-z])
GOARGS=GOARCH=amd64 CGO_ENABLED=0
GOBUILD=$(GOARGS) $(GOBIN) build -ldflags="$(RELEASE)"

build: unix
darwin:
	GOOS=darwin $(GOBUILD) ./cmd/code2mysql/
	GOOS=darwin $(GOBUILD) ./cmd/table2file/
mac: darwin
macos: darwin
linux:
	GOOS=linux $(GOBUILD) ./cmd/code2mysql/
	GOOS=linux $(GOBUILD) ./cmd/table2file/
unix:
	GOOS=$(GOOS) $(GOBUILD) ./cmd/code2mysql/
	GOOS=$(GOOS) $(GOBUILD) ./cmd/table2file/
upx: linux
	$(UPXBIN) code2mysql
	$(UPXBIN) table2file
upxx: linux
	$(UPXBIN) --ultra-brute code2mysql
	$(UPXBIN) --ultra-brute table2file
win: windows
windows:
	GOOS=windows $(GOBUILD) ./cmd/code2mysql/
	GOOS=windows $(GOBUILD) ./cmd/table2file/
