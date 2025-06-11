package hasher

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
)

// FileHash 结构体用于存储文件的路径、哈希值和大小
type FileHash struct {
	Path string
	Hash string
	Size int64
}

// FormatFileSize 格式化文件大小为可读的字符串
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Worker 函数用于计算文件的哈希值
func Worker(files <-chan string, hashes chan<- FileHash, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range files {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("打开文件 %s 出错: %v\n", path, err)
			continue
		}

		info, err := file.Stat()
		if err != nil {
			fmt.Printf("获取文件 %s 信息出错: %v\n", path, err)
			file.Close()
			continue
		}

		h := md5.New()
		if _, err := io.Copy(h, file); err != nil {
			fmt.Printf("计算文件 %s 哈希值出错: %v\n", path, err)
			file.Close()
			continue
		}
		file.Close()

		hashSum := hex.EncodeToString(h.Sum(nil))
		hashes <- FileHash{Path: path, Hash: hashSum, Size: info.Size()}
	}
}
