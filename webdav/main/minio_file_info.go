package main

import (
	"github.com/minio/minio-go/v7"
	"os"
	"time"
)

// minioFileInfo 实现了 os.FileInfo 接口
type minioFileInfo struct {
	objInfo minio.ObjectInfo
}

// Name 返回文件名
func (fi *minioFileInfo) Name() string {
	return fi.objInfo.Key
}

// Size 返回文件大小
func (fi *minioFileInfo) Size() int64 {
	return fi.objInfo.Size
}

// Mode 返回文件权限和模式
func (fi *minioFileInfo) Mode() os.FileMode {
	return 0644 // 示例中使用默认权限
}

// ModTime 返回文件修改时间
func (fi *minioFileInfo) ModTime() time.Time {
	return fi.objInfo.LastModified
}

// IsDir 返回是否是目录
func (fi *minioFileInfo) IsDir() bool {
	return false // 示例中假设文件不是目录
}

// Sys 返回底层数据源（可以返回 nil）
func (fi *minioFileInfo) Sys() interface{} {
	return nil
}
