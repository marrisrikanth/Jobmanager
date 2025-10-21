# ================================
# Job Manager Makefile
# ================================

APP_NAME      := jobmanager
PKG           := ./...
BUILD_DIR     := build
VERSION       := v0.0.0
COMMIT        := unknown
DATE          := 
LDFLAGS       := -s -w -X 'main.Version=' -X 'main.Commit=' -X 'main.BuildDate='

# default target
.PHONY: all
all: clean build

# ================================
# Build Targets
# ================================

## 🧱 Normal build (optimized)
.PHONY: build
build:
@echo "➡️  Building  (Version: , Commit: )"
CGO_ENABLED=0 go build -trimpath -ldflags "" -o / .

## 🧱 Development build (with debug symbols)
.PHONY: debug
debug:
go build -o /-debug .

## 🧱 Cross-compile for Linux, Windows, macOS
.PHONY: cross
cross:
@echo "➡️  Cross-compiling for multiple OS/ARCH..."
mkdir -p 
GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "" -o /-linux-amd64 .
GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "" -o /-darwin-amd64 .
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "" -o /-windows-amd64.exe .

## 🧱 Run the binary locally
.PHONY: run
run: build
.//

## 🧹 Clean build artifacts
.PHONY: clean
clean:
rm -rf 

# ================================
# Compression (Optional)
# ================================

## 🧩 Compress binary with UPX
.PHONY: compress
compress: build
@echo "➡️  Compressing binary with UPX..."
upx --brute / || echo "⚠️  UPX not installed or failed (skipping)"

# ================================
# Docker
# ================================

## 🐳 Build minimal Docker image
.PHONY: docker
docker:
@echo "➡️  Building Docker image..."
docker build -t : .

# ================================
# Test and Benchmark
# ================================

## 🧪 Run tests
.PHONY: test
test:
go test -v 

## ⚡ Run benchmarks
.PHONY: bench
bench:
go test -bench=. 

# ================================
# Release (optimized + compressed)
# ================================

## 🚀 Release build
.PHONY: release
release: clean build compress
@echo "✅ Release binary ready at /"


