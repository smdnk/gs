package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

var accessKey = "BFHpp7CcxnBKv58b4XSa"
var accessSecret = "yqysI79US6QpuTNlO1Y6bk7DSu4hlpmK9DKX7jR7"
var endPoint = "192.168.0.6:9000" // minio地址,不能加http

var minioClient = initMinioClient()

func initMinioClient() *minio.Client {
	// 初始化minio客户端
	minioClient, err := minio.New(endPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, accessSecret, ""),
		Secure: false, // 是否使用https进行通信
	})
	if err != nil {
		log.Fatal("minio client create fail, err %+v", err)
	}

	return minioClient
}

func main() {
	//bucketObjectList("test")
	objectInfo("test", "smdnk/wallhaven-3z88ld_3440x1440.png")

}

// creatBucket 创建 bucket
func creatBucket(bucketName string) {
	err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("Successfully created %s.", bucketName)
}

// bucketInfo 获取bucket信息
func bucketInfo(bucketName string) {
	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	var bucket minio.BucketInfo
	for _, v := range buckets {
		if v.Name == bucketName {
			bucket = v
		}
	}

	// 打印存储桶的基本信息
	log.Printf("Bucket name: %s\n", bucket.Name)
	//log.Printf("Creation time: %s\n", bucketInfo.CreationDate)
	//log.Printf("Region: %s\n", bucketInfo.Region)
	//log.Printf("Object count: %d\n", bucketInfo.ObjectCount)
	//log.Printf("Size: %d bytes\n", bucketInfo.Size)
}

// 列出所有bucket
func bucketList() {
	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, bucket := range buckets {
		log.Println(bucket)
	}
}

// 判断bucket是否存在
func bucketIsHas(bucketName string) {
	found, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if found {
		log.Println("Bucket found")
	}
}

// 删除bucket
func bucketDel(bucketName string) {
	err := minioClient.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("Bucket remove: %s", bucketName)
}

// 从本地读入文件并上传
func objectUpLocal(fileName string, bucketName string, objectName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	/*
	   minio.UploadInfo
	   	info.ETag | string | The ETag of the new object
	   	info.VersionID | string | The version identifyer of the new object
	*/
	UserMetadata := map[string]string{
		"origin_name": "etcd-v3.4.18-linux-amd64.tar.gz",
	}
	objectSize := fileStat.Size() // objectSize可设置为-1，表示不确定文件大小，但是-1会预分配比较大的内存。
	uploadInfo, err := minioClient.PutObject(context.Background(), bucketName, objectName, file, objectSize,
		minio.PutObjectOptions{ContentType: "application/octet-stream", UserMetadata: UserMetadata})
	/*
		minio.PutObjectOptions用得比较多的是这几个参数，注意一下：
		ContentType：string，用于设置下载时Response的header里的ContentType，如"application/octet-stream"，"image/jpg"；
		WebsiteRedirectLocation：string，重定向URL，比如客户端下载的是a，此时可以通过该参数的url来重定向url下载b，这样就可以实现不需客户端修改下载地址的情况下实现下载其他资源。
		UserMetadata：map[string]string，用于存储一些自定义的额外信息，比如你想存储object所关联的uid，就把"uid":123444 设置进去，完成文件与信息的关联。
	*/
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)
}

// 从网络IO流中上传文件 （以gin框架为例）
/*func objectUpNet() {
	file, err := c.FormFile("file") // 从网络IO中获取文件流
	if err != nil {
		return errors.WithCode(errcode.ErrUnknown, err.Error()), 0
	}

	reader, err := file.Open()
	if err != nil {
		log.Println(err)
		return
	}
	defer reader.Close()

	objectSize := -1 // 表示object大小未知
	_, err := client.PutObject(context.Background(), bucketName, objectName, reader, objectSize,
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println(err)
		return
	}
}*/

// 列出指定bucket下的所有obejct
func bucketObjectList(bucketName string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注意，ListObjects返回值是个channel，通过迭代来获取所有object
	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    "", // 通过该参数过滤以Prefix作为object key前缀的object
		Recursive: true,
	})
	/*
		minio.ObjectInfo
		objectInfo.Key | string | Name of the object
		objectInfo.Size | int64 | Size of the object
		objectInfo.ETag | string | MD5 checksum of the object
		objectInfo.LastModified | time.Time | Time when object was last modified
	*/
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		log.Println(object.Key)
	}
}

// 获取object的基本信息
func objectInfo(bucketName string, objectName string) {
	/*
	   minio.ObjectInfo

	   	objectInfo.Key | string | Name of the object
	   	objectInfo.Size | int64 | Size of the object
	   	objectInfo.ETag | string | MD5 checksum of the object
	   	objectInfo.LastModified | time.Time | Time when object was last modified
	*/
	objInfo, err := minioClient.StatObject(context.Background(), bucketName, objectName,
		minio.StatObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("stat:%+v", objInfo)

}

// 删除object
func objectDel(bucketName string, objectName string) {
	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 从oss复制一份object
func copyObject(bucketName string, objectName string) {
	// 源object
	srcOpts := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: objectName,
	}

	// 目标object
	dstOpts := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: objectName + "-copy",
	}

	uploadInfo, err := minioClient.CopyObject(context.Background(), dstOpts, srcOpts)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully copied object:", uploadInfo)
}
