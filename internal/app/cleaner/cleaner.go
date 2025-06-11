package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"clean_repeat_files/internal/app/hasher"
)

// FindAndCleanDuplicates 查找并清理重复文件
func FindAndCleanDuplicates(dirPath string) (int, int64, error) {
	fileHashes := sync.Map{}
	files := make(chan string)
	hashes := make(chan hasher.FileHash)
	doneCollecting := make(chan bool)
	var duplicateCount int
	var estimatedCleanupSize int64

	// Start workers
	var wg sync.WaitGroup
	numWorkers := 5 // 可以根据CPU核心数调整
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go hasher.Worker(files, hashes, &wg)
	}

	// Collect results and count duplicates
	go func() {
		for fh := range hashes {
			if _, loaded := fileHashes.LoadOrStore(fh.Hash, fh.Path); loaded {
				duplicateCount++
				estimatedCleanupSize += fh.Size
			}
		}
		close(doneCollecting)
	}()

	// Walk the directory and send files to workers
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files <- path
		}
		return nil
	})

	if err != nil {
		return 0, 0, fmt.Errorf("遍历目录出错: %v", err)
	}

	close(files)
	wg.Wait()
	close(hashes)
	<-doneCollecting

	return duplicateCount, estimatedCleanupSize, nil
}

// PerformCleanup 执行清理操作
func PerformCleanup(dirPath string) (int64, error) {
	fileHashes := sync.Map{}
	files := make(chan string)
	hashes := make(chan hasher.FileHash)
	doneCleanup := make(chan bool)
	var actualCleanupSize int64

	// Restart workers for cleanup
	var wg sync.WaitGroup
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go hasher.Worker(files, hashes, &wg)
	}

	// Collect results and perform cleanup
	go func() {
		for fh := range hashes {
			if existingPath, loaded := fileHashes.LoadOrStore(fh.Hash, fh.Path); loaded {
				// 发现重复文件，删除它
				// 为了安全，我们只删除重复文件，保留原始文件
				err := os.Remove(fh.Path)
				if err != nil {
					fmt.Printf("删除文件 %s 出错: %v\n", fh.Path, err)
				} else {
					actualCleanupSize += fh.Size
					fmt.Printf("已删除重复文件: %s (原文件: %s)\n", fh.Path, existingPath)
				}
			} else {
				// 存储原始文件路径，以便后续比较
				fileHashes.Store(fh.Hash, fh.Path)
			}
		}
		close(doneCleanup)
	}()

	// Walk the directory again and send files to workers for cleanup
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files <- path
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("遍历目录出错: %v", err)
	}

	close(files)
	wg.Wait()
	close(hashes)
	<-doneCleanup

	return actualCleanupSize, nil
}
