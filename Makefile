BINNAME=dbtool
RELEASE=-s -w
UPXBIN=/usr/local/bin/upx
GOBIN=/usr/local/bin/go
GOOS=$(shell uname -s | tr [A-Z] [a-z])
GOARGS=GOARCH=amd64 CGO_ENABLED=0
GOBUILD=$(GOARGS) $(GOBIN) build -ldflags="$(RELEASE)"

build: unix
darwin:
	GOOS=darwin $(GOBUILD) -o $(BINNAME) .
mac: darwin
macos: darwin
linux:
	GOOS=linux $(GOBUILD) -o $(BINNAME) .
unix:
	GOOS=$(GOOS) $(GOBUILD) -o $(BINNAME) .
upx: linux
	$(UPXBIN) $(BINNAME)
upxx: linux
	$(UPXBIN) --ultra-brute $(BINNAME)
win: windows
windows:
	GOOS=windows $(GOBUILD) -o $(BINNAME).exe .
