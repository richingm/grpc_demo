package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

func GetRealFilePath(cwd, exitDir, relativePath string, parentPath string, fileName string) string {
	// 打印当前工作目录，方便调试
	fmt.Println("Current working directory:", cwd)

	// 获取项目根目录的路径
	projectRoot := getProjectRoot(cwd, exitDir)

	// 需要加载的证书文件路径（相对路径）
	filePath := filepath.Join(projectRoot, relativePath, parentPath, fileName)
	return filePath
}

// getProjectRoot 根据当前工作目录推测项目根目录
func getProjectRoot(cwd string, exitDir string) string {
	if strings.Contains(cwd, exitDir) {
		// 查找 "tls_demo" 在路径中的位置
		index := strings.Index(cwd, exitDir)

		if index == -1 {
			fmt.Println("No 'tls_demo' found in the path.")
			return cwd
		}

		// 截取 "tls_demo" 之前的部分
		cwd = cwd[:index+len(exitDir)]
		return cwd
	}
	return fmt.Sprintf("%s/%s", cwd, exitDir)
}
