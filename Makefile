# xcap Makefile

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: all build clean

# 默认目标
all: build

# 构建当前平台
build:
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/xcap ./cmd/xcap/

# 清理构建产物
clean:
	@rm -rf bin/
	@rm -rf output/
