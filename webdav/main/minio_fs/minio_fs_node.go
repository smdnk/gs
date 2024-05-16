package minio_fs

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/minio/minio-go/v7"
	"gs/webdav"
	"net/http"
	"os"
	"sync"
	"time"
)

type minioFSNode struct {
	children   map[string]*minioFSNode
	client     *minio.Client
	bucketName string
	nodeFlg    int
	mu         sync.Mutex
	data       []byte
	mode       os.FileMode
	modTime    time.Time
	deadProps  map[xml.Name]webdav.Property
}

func (n *minioFSNode) stat(name string) *minioFileInfo {
	n.mu.Lock()
	defer n.mu.Unlock()

	size := int64(len(n.data))
	if len(name) > 3 && n.mode == os.ModeDir {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// 注意，ListObjects返回值是个channel，通过迭代来获取所有object
		objectCh := n.client.ListObjects(ctx, name, minio.ListObjectsOptions{
			Prefix:    "", // 通过该参数过滤以Prefix作为object key前缀的object
			Recursive: true,
		})

		for object := range objectCh {
			if object.Err != nil {
				panic(object.Err)
			}
			size += object.Size
		}
	}

	if n.mode != os.ModeDir {
		objInfo, err := n.client.StatObject(context.Background(), n.bucketName, name,
			minio.StatObjectOptions{})
		if err != nil {
			fmt.Println(err)
		}
		size = objInfo.Size
	}

	return &minioFileInfo{
		name:    name,
		size:    size,
		mode:    n.mode,
		modTime: n.modTime,
	}
}

func (n *minioFSNode) DeadProps() (map[xml.Name]webdav.Property, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if len(n.deadProps) == 0 {
		return nil, nil
	}
	ret := make(map[xml.Name]webdav.Property, len(n.deadProps))
	for k, v := range n.deadProps {
		ret[k] = v
	}
	return ret, nil
}

func (n *minioFSNode) Patch(patches []webdav.Proppatch) ([]webdav.Propstat, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	pstat := webdav.Propstat{Status: http.StatusOK}
	for _, patch := range patches {
		for _, p := range patch.Props {
			pstat.Props = append(pstat.Props, webdav.Property{XMLName: p.XMLName})
			if patch.Remove {
				delete(n.deadProps, p.XMLName)
				continue
			}
			if n.deadProps == nil {
				n.deadProps = map[xml.Name]webdav.Property{}
			}
			n.deadProps[p.XMLName] = p
		}
	}
	return []webdav.Propstat{pstat}, nil
}
