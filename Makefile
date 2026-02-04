# xcap Makefile

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

# macOS: 抑制 duplicate libraries 警告
ifeq ($(shell uname),Darwin)
    export CGO_LDFLAGS := -Wl,-no_warn_duplicate_libraries
endif

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
