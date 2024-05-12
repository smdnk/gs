package main

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"os"
)

// minioFile 实现了 webdav.File 接口
type minioFile struct {
	client *minio.Client
	name   string
}

func (f *minioFile) Close() error {
	// 在MinIO中关闭文件的操作
	fmt.Println("minio Close")
	return nil
}

func (f *minioFile) Read(p []byte) (int, error) {
	// 从MinIO中读取文件的操作
	fmt.Println("minio Read")
	return 0, nil
}

func (f *minioFile) Seek(offset int64, whence int) (int64, error) {
	// 在MinIO中移动文件指针的操作
	fmt.Println("minio Seek")
	return 0, nil
}

func (f *minioFile) Write(p []byte) (int, error) {
	// 向MinIO中写入文件的操作
	fmt.Println("minio Write")
	return 0, nil
}

func (f *minioFile) Readdir(count int) ([]os.FileInfo, error) {
	// 读取MinIO目录中文件列表的操作
	fmt.Println("minio Readdir")
	return nil, nil
}

func (f *minioFile) Stat() (os.FileInfo, error) {
	// 获取MinIO文件信息的操作
	fmt.Println("minio Stat")
	return nil, nil
}
