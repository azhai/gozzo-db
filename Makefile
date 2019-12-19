RELEASE=-s -w
UPXBIN=/usr/local/bin/upx
GOBIN=/usr/local/bin/go
GOOS=$(shell uname -s | tr [A-Z] [a-z])
GOARGS=GOARCH=amd64 CGO_ENABLED=0
GOBUILD=$(GOARGS) $(GOBIN) build -ldflags="$(RELEASE)"
TARGETS=$(shell ls -1 -I '*.*' ./cmd/ | grep -v tmp | xargs)

.PHONY: clean common
clean:
    rm -f $(TARGETS)
common:
    $(foreach dst, $(TARGETS), GOOS=$(GOOS) $(GOBUILD) ./cmd/$(dst)/)
darwin:
    $(foreach dst, $(TARGETS), GOOS=darwin $(GOBUILD) ./cmd/$(dst)/)
mac: darwin
macos: darwin
linux:
    $(foreach dst, $(TARGETS), GOOS=linux $(GOBUILD) ./cmd/$(dst)/)
upx: linux
	$(UPXBIN) $(TARGETS)
upxx: linux
	$(UPXBIN) --ultra-brute $(TARGETS)
win: windows
windows:
    $(foreach dst, $(TARGETS), GOOS=windows $(GOBUILD) ./cmd/$(dst)/)
