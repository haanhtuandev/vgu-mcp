BINARY  := vgu-mcp
CMD     := ./cmd/vgu-mcp
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

# 5 distribution targets: windows/arm64 intentionally excluded
TARGETS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64

.PHONY: build run tidy test install dist clean

build:
	go build $(LDFLAGS) -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

tidy:
	go mod tidy

test:
	go test ./...

install: build
	install -m 755 $(BINARY) /usr/local/bin/$(BINARY)
	@echo "✓ Installed to /usr/local/bin/$(BINARY)"

dist: clean
	@mkdir -p dist
	@for target in $(TARGETS); do \
		os=$$(echo $$target | cut -d/ -f1); \
		arch=$$(echo $$target | cut -d/ -f2); \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		outdir="dist/$(BINARY)-$(VERSION)-$$os-$$arch"; \
		mkdir -p $$outdir; \
		GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) \
			-o $$outdir/$(BINARY)$$ext $(CMD); \
		echo "✓ $$outdir/$(BINARY)$$ext"; \
	done

clean:
	rm -rf dist $(BINARY)
