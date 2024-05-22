package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"gs/go_webdav"
	"log"
	"os"
	"slices"
	"sync"
)

type MinioFS struct {
	client             *minio.Client
	mu                 sync.Mutex
	bucketList         map[string]*Bucket
	bucketNameList     []string
	currentBucket      string
	currentObjectNames []string
}
type Bucket struct {
	bucketName string
	objects    map[string]*minio.ObjectInfo
	isFile     bool
}

func NewMinioFS(client *minio.Client) go_webdav.FileSystem {
	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chl := make(map[string]*Bucket, 10)
	var bucketNames []string
	for _, bucket := range buckets {
		bucketNames = append(bucketNames, bucket.Name)
		chl[bucket.Name] = &Bucket{
			bucketName: bucket.Name,
			isFile:     true,
			objects:    make(map[string]*minio.ObjectInfo, 10),
		}
		// 注意，ListObjects返回值是个channel，通过迭代来获取所有object
		objectCh := client.ListObjects(ctx, bucket.Name, minio.ListObjectsOptions{
			Prefix:    "", // 通过该参数过滤以Prefix作为object key前缀的object
			Recursive: true,
		})
		for object := range objectCh {
			if object.Err != nil {
				fmt.Println(object.Err)
			}
			chl[bucket.Name].objects[object.Key] = &object
		}
	}

	return &MinioFS{
		bucketList:     chl,
		bucketNameList: bucketNames,
		client:         client,
		currentBucket:  "/",
	}
}

func (fs *MinioFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {

	return nil
}

func (fs *MinioFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (go_webdav.File, error) {
	return nil, nil
}

func (fs *MinioFS) RemoveAll(ctx context.Context, name string) error {

	return nil
}

func (fs *MinioFS) Rename(ctx context.Context, oldName, newName string) error {

	return nil
}

func (fs *MinioFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {

	return nil, os.ErrNotExist
}

func (fs *MinioFS) CurrentFileList(ctx context.Context, name string) ([]string, error) {
	// 如果是根目录 返回bucket列表
	if name == "/" || name == "" {
		return fs.bucketNameList, nil
	}
	// 如果是bucket名字 返回object列表
	if contains := slices.Contains(fs.bucketNameList, name); contains {
		currentBucket := fs.bucketList[fs.currentBucket]
		objects := currentBucket.objects

		var objectNames []string
		for k, _ := range objects {
			objectNames = append(objectNames, k)
		}
		fs.currentObjectNames = objectNames

		log.Println(name)

		return objectNames, nil
	}

	// 如果是对象名字 返回对象信息
	if contains := slices.Contains(fs.currentObjectNames, name); contains {

	}

	return nil, nil

}

func (fs *MinioFS) SetCurrentBucket(ctx context.Context, bucketName string) {
	contains := slices.Contains(fs.bucketNameList, bucketName)
	// todo 如果某个目录和bucket名字相同 这里就会有bug
	if contains {
		fs.currentBucket = bucketName
	}
}
