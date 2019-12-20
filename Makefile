RELEASE=-s -w
UPXBIN=/usr/local/bin/upx
GOBIN=/usr/local/bin/go
GOOS=$(shell uname -s | tr [A-Z] [a-z])
GOARGS=GOARCH=amd64 CGO_ENABLED=0
APPS=$(shell ls -1 ./cmd/ | grep -v "\." | grep -v tmp | xargs)

.PHONY: all
all: clean build
build:
	@ for dst in $(APPS); \
	do \
		GOOS=$(GOOS) $(GOARGS) $(GOBIN) build -ldflags="$(RELEASE)" ./cmd/$$dst/; \
	done
	@echo Compile $(APPS)
	@echo Build success.
clean:
	rm -f $(APPS)
	@echo Clean all.
upx: build
	$(UPXBIN) $(APPS)
upxx: build
	$(UPXBIN) --ultra-brute $(APPS)

