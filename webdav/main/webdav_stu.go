package main

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gs/webdav"
	"gs/webdav/main/minio_fs"
	"io/fs"
	"log"
	"net/http"
)

func main() {
	// 初始化MinIO客户端
	minioClient, err := minio.New("minioapi.smdnk.cn:80", &minio.Options{
		Creds:  credentials.NewStaticV4("BFHpp7CcxnBKv58b4XSa", "yqysI79US6QpuTNlO1Y6bk7DSu4hlpmK9DKX7jR7", ""),
		Secure: false, // 默认情况下，MinIO没有启用HTTPS，可以根据需要修改
	})
	if err != nil {
		log.Fatalln(err)
	}

	newMinioFS := minio_fs.NewMinioFS(minioClient)
	//dirFs := minio_dir.MinioDir{
	//	DirName: "./",
	//}

	buckets, err := minioClient.ListBuckets(context.Background())
	for _, bucket := range buckets {
		newMinioFS.Mkdir(context.Background(), bucket.Name, fs.ModeDir)
	}

	// 创建一个新的WebDAV处理器
	handler := &webdav.Handler{
		Prefix: "/",
		//FileSystem: dirFs,
		FileSystem: newMinioFS,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			log.Printf("%s %s: %v", r.Method, r.URL.Path, err)
		},
	}

	// 设置HTTP服务器
	http.Handle("/", handler)

	// 启动HTTP服务器
	log.Fatal(http.ListenAndServe(":8080", nil))
}
