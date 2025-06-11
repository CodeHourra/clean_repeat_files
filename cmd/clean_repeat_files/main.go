package main

import (
	"flag"
	"fmt"
	"os"

	"clean_repeat_files/internal/app/cleaner"
	"clean_repeat_files/internal/app/hasher"
	"clean_repeat_files/internal/app/ui"
)

func main() {
	// Parse command-line arguments
	var dirPath string
	flag.StringVar(&dirPath, "path", "", "需要扫描的目录路径")
	flag.Parse()

	// 交互式选择目录
	if dirPath == "" {
		selectedPath, err := ui.SelectDirectory()
		if err != nil {
			fmt.Printf("选择目录失败: %v\n", err)
			os.Exit(1)
		}
		dirPath = selectedPath
	}

	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Printf("目录不存在: %s\n", dirPath)
		os.Exit(1)
	}

	// 添加loading效果
	doneLoading := make(chan bool)
	go ui.ShowLoading(doneLoading)

	// 查找重复文件
	duplicateCount, estimatedCleanupSize, err := cleaner.FindAndCleanDuplicates(dirPath)
	if err != nil {
		fmt.Printf("查找重复文件出错: %v\n", err)
		doneLoading <- true
		os.Exit(1)
	}

	doneLoading <- true

	// 显示重复文件数量并询问是否继续清理
	fmt.Printf("发现 %d 个重复文件。预计可清理空间: %s。是否继续清理？(y/n): ", duplicateCount, hasher.FormatFileSize(estimatedCleanupSize))
	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		return
	}

	if response == "y" || response == "Y" {
		actualCleanupSize, err := cleaner.PerformCleanup(dirPath)
		if err != nil {
			fmt.Printf("执行清理出错: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("清理完成。实际清理空间: %s\n", hasher.FormatFileSize(actualCleanupSize))
	} else {
		fmt.Println("清理操作已取消。")
	}
}
