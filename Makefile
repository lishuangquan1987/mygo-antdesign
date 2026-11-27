# mygo-antdesign 开发命令（可选；不依赖 Make 时可直接执行对应 go 命令）。
SHELL := bash
GO ?= go
VERSION ?= dev
BIN := bin/gallery
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: all help run build test vet fmt check baselines clean

all: check

help:
	@echo "make run        运行 Gallery 示例应用（go run ./gallery）"
	@echo "make build      构建 Gallery 到 $(BIN)"
	@echo "make test       运行全部测试（go test -count=1 ./...）"
	@echo "make vet        go vet ./..."
	@echo "make fmt        格式化代码（gofmt -w .）"
	@echo "make check      fmt 检查 + vet + test + 样式铁律"
	@echo "make baselines  重新生成截图基线（Gallery 级 + 逐示例）"
	@echo "make clean      删除构建产物 bin/"

run:
	$(GO) run ./gallery

build:
	mkdir -p bin
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) ./gallery
	@echo "built $(BIN)"

test:
	$(GO) test -count=1 ./...

vet:
	$(GO) vet ./...

fmt:
	gofmt -w .

check:
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then echo "未格式化：$$out"; exit 1; fi
	$(GO) vet ./...
	$(GO) test -count=1 ./...
	$(GO) run ./internal/check

baselines:
	$(GO) test ./gallery ./samples -count=1 -run 'TestGalleryBaselines|TestControlScreenshots'

clean:
	rm -rf bin
