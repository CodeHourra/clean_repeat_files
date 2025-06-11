# 定义可执行文件名
BINARY_NAME := clean_repeat_files
# Go编译参数
GO_BUILD_FLAGS := -v
# 安装路径
INSTALL_PATH := /usr/local/bin

# 默认目标：编译程序
build:
	go build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/clean_repeat_files

# 清理生成的可执行文件
clean:
	rm -f $(BINARY_NAME)*

# 安装程序到系统路径
install: build
	cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

# 使用goreleaser进行发布
release:
	goreleaser release --clean

# 使用goreleaser进行快照构建（不创建tag）
snapshot:
	goreleaser release --snapshot --clean

.PHONY: build clean install release snapshot