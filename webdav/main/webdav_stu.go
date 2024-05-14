package main

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gs/webdav"
	"log"
	"net/http"
)

func main() {
	// 初始化MinIO客户端
	minioClient, err := minio.New("192.168.0.6:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("BFHpp7CcxnBKv58b4XSa", "yqysI79US6QpuTNlO1Y6bk7DSu4hlpmK9DKX7jR7", ""),
		Secure: false, // 默认情况下，MinIO没有启用HTTPS，可以根据需要修改
	})
	if err != nil {
		log.Fatalln(err)
	}

	fs := &minioFS{
		client: minioClient,
	}

	// 创建一个新的WebDAV处理器
	handler := &webdav.Handler{
		Prefix:     "/",
		FileSystem: fs,
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
