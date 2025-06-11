package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
)

// presetDirs 定义了预设目录，方便用户选择
var presetDirs = map[string]string{
	"微信":      "~/Library/Containers/com.tencent.xinWeChat/Data/Library/Application Support/com.tencent.xinWeChat",
	"企业微信":    "~/Library/Containers/com.tencent.WeWorkMac/Data/Library/Application Support/com.tencent.WeWorkMac",
	"其他目录...": "",
}

// SelectDirectory 交互式选择目录
func SelectDirectory() (string, error) {
	items := make([]string, 0, len(presetDirs))
	for k := range presetDirs {
		items = append(items, k)
	}

	prompt := promptui.Select{
		Label: "请选择要扫描的目录",
		Items: items,
		Size:  5,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	if result == "其他目录..." {
		return promptForCustomDir()
	}

	// 展开预设目录中的~为实际路径
	dir, err := homedir.Expand(presetDirs[result])
	if err != nil {
		return "", fmt.Errorf("解析路径失败: %v", err)
	}

	return dir, nil
}

// promptForCustomDir 提示用户输入自定义目录
func promptForCustomDir() (string, error) {
	prompt := promptui.Prompt{
		Label:   "请输入目录路径",
		Default: "~/Downloads",
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("目录不能为空")
			}
			return nil
		},
	}

	dir, err := prompt.Run()
	if err != nil {
		return "", err
	}

	// 展开路径中的~为实际路径
	expandedDir, err := homedir.Expand(dir)
	if err != nil {
		return "", fmt.Errorf("解析路径失败: %v", err)
	}

	return expandedDir, nil
}

// ShowLoading 显示加载动画
func ShowLoading(done chan bool) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Printf("\r%s\r", strings.Repeat(" ", 20)) // 清除loading显示
			return
		default:
			fmt.Printf("\r正在统计重复文件 %s", frames[i])
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(frames)
		}
	}
}
