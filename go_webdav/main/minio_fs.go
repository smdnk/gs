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
	buckets            map[string]*go_webdav.MinioFileInfo
	bucketNameList     []string
	currentBucket      string
	currentObjects     []*go_webdav.MinioFileInfo
	currentObjectNames []string
}

func NewMinioFS(client *minio.Client) go_webdav.FileSystem {
	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chl := make(map[string]*go_webdav.MinioFileInfo, 10)
	var bucketNames []string
	for _, bucket := range buckets {
		bucketNames = append(bucketNames, bucket.Name)
		chl[bucket.Name] = &go_webdav.MinioFileInfo{
			BucketName: bucket.Name,
			Objects:    make(map[string]*minio.ObjectInfo, 10),
			Siz:        1, // todo
			Mod:        os.ModeDir,
			ModTim:     bucket.CreationDate,
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
			chl[bucket.Name].Objects[object.Key] = &object
		}
	}

	return &MinioFS{
		buckets:        chl,
		bucketNameList: bucketNames,
		client:         client,
		currentBucket:  "/",
	}
}

func (fs *MinioFS) CurrentFileList(ctx context.Context, name string) ([]*go_webdav.MinioFileInfo, error) {
	// 如果是根目录 返回bucket列表
	if name == "/" || name == "" {
		var fileInfo []*go_webdav.MinioFileInfo
		for _, bucketInfo := range fs.buckets {
			fileInfo = append(fileInfo, bucketInfo)
		}
		return fileInfo, nil
	}
	// 如果是bucket名字 返回object列表
	if contains := slices.Contains(fs.bucketNameList, name); contains {
		currentBucket := fs.buckets[fs.currentBucket]
		objects := currentBucket.Objects

		var objectNames []string
		var fileInfoList []*go_webdav.MinioFileInfo
		for objName, objInfo := range objects {
			fileInfo := &go_webdav.MinioFileInfo{
				BucketName: objName,
				Objects:    make(map[string]*minio.ObjectInfo, 1),
				Siz:        objInfo.Size,
				Mod:        os.ModeDir,
				ModTim:     objInfo.LastModified,
			}
			fileInfoList = append(fileInfoList, fileInfo)
			objectNames = append(objectNames, objName)
		}
		fs.currentObjects = fileInfoList
		fs.currentObjectNames = objectNames

		log.Println(name)

		return fileInfoList, nil
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
