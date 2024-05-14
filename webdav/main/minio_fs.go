package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"gs/webdav"
	"os"
)

type minioFS struct {
	client *minio.Client
}

func (fs minioFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {

	fmt.Println(name, perm, "Mkdir")
	return nil
}

func (fs minioFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	// 在MinIO中打开文件的操作
	fmt.Println(name, flag, perm, "OpenFile")
	//f, err := os.OpenFile(name, flag, perm)
	return &minioFile{
		client: fs.client,
		name:   name,
	}, nil
}

func (fs minioFS) RemoveAll(ctx context.Context, name string) error {

	fmt.Println(name, "RemoveAll")
	return nil
}

func (fs minioFS) Rename(ctx context.Context, oldName, newName string) error {

	fmt.Println(oldName, newName, "Rename")
	return nil
}

func (fs minioFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	fmt.Println(name, "Stat")
	objInfo, err := fs.client.StatObject(ctx, "test", name, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &minioFileInfo{objInfo}, nil
}
