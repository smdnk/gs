package minio_fs

import (
	"github.com/minio/minio-go/v7"
	"os"
	"time"
)

// minioFileInfo 实现了 os.FileInfo 接口
type minioFileInfo struct {
	objInfo minio.ObjectInfo
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (f *minioFileInfo) Name() string       { return f.name }
func (f *minioFileInfo) Size() int64        { return f.size }
func (f *minioFileInfo) Mode() os.FileMode  { return f.mode }
func (f *minioFileInfo) ModTime() time.Time { return f.modTime }
func (f *minioFileInfo) IsDir() bool        { return f.mode.IsDir() }
func (f *minioFileInfo) Sys() interface{}   { return nil }
