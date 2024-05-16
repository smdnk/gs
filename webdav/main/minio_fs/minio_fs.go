package minio_fs

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"gs/webdav"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

type MinioFS struct {
	client *minio.Client
	mu     sync.Mutex
	root   minioFSNode
}

func NewMinioFS(client *minio.Client) webdav.FileSystem {
	return &MinioFS{
		root: minioFSNode{
			client:   client,
			nodeFlg:  0,
			children: make(map[string]*minioFSNode),
			mode:     0660 | os.ModeDir,
			modTime:  time.Now(),
		},
		client: client,
	}
}

// The frag argument will be empty only if dir is the root node and the walk
// ends at that root node.
func (fs *MinioFS) walk(op, fullname string, f func(dir *minioFSNode, frag string, final bool) error) error {
	original := fullname
	fullname = slashClean(fullname)

	// Strip any leading "/"s to make fullname a relative path, as the walk
	// starts at fs.root.
	if fullname[0] == '/' {
		fullname = fullname[1:]
	}
	dir := &fs.root

	for {
		frag, remaining := fullname, ""
		i := strings.IndexRune(fullname, '/')
		final := i < 0
		if !final {
			frag, remaining = fullname[:i], fullname[i+1:]
		}
		if frag == "" && dir != &fs.root {
			panic("webdav: empty path fragment for a clean path")
		}
		if err := f(dir, frag, final); err != nil {
			return &os.PathError{
				Op:   op,
				Path: original,
				Err:  err,
			}
		}
		if final {
			break
		}
		child := dir.children[frag]
		if child == nil {
			return &os.PathError{
				Op:   op,
				Path: original,
				Err:  os.ErrNotExist,
			}
		}
		if !child.mode.IsDir() {
			return &os.PathError{
				Op:   op,
				Path: original,
				Err:  os.ErrInvalid,
			}
		}
		dir, fullname = child, remaining
	}
	return nil
}

func (fs *MinioFS) find(op, fullname string) (parent *minioFSNode, frag string, err error) {
	err = fs.walk(op, fullname, func(parent0 *minioFSNode, frag0 string, final bool) error {
		if !final {
			return nil
		}
		if frag0 != "" {
			parent, frag = parent0, frag0
		}
		return nil
	})
	return parent, frag, err
}

func (fs *MinioFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dir, frag, err := fs.find("mkdir", name)
	if err != nil {
		return err
	}
	if dir == nil {
		// We can't create the root.
		return os.ErrInvalid
	}
	if _, ok := dir.children[frag]; ok {
		return os.ErrExist
	}
	dir.children[frag] = &minioFSNode{
		client:     fs.client,
		bucketName: name,
		children:   make(map[string]*minioFSNode),
		mode:       perm.Perm() | os.ModeDir,
		modTime:    time.Now(),
	}
	// minio 新建文件夹
	return nil
}

func (fs *MinioFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dir, frag, err := fs.find("open", name)
	if err != nil {
		return nil, err
	}
	var n *minioFSNode
	if dir == nil {
		// We're opening the root.
		if runtime.GOOS == "zos" {
			if flag&os.O_WRONLY != 0 {
				return nil, os.ErrPermission
			}
		} else {
			if flag&(os.O_WRONLY|os.O_RDWR) != 0 {
				return nil, os.ErrPermission
			}
		}
		n, frag = &fs.root, "/"

	} else {
		n = dir.children[frag]

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// 注意，ListObjects返回值是个channel，通过迭代来获取所有object
		objectCh := dir.client.ListObjects(ctx, frag, minio.ListObjectsOptions{
			Prefix:    "", // 通过该参数过滤以Prefix作为object key前缀的object
			Recursive: true,
		})
		for object := range objectCh {
			if object.Err != nil {
				fmt.Println(object.Err)
			}
			n.children[object.Key] = &minioFSNode{
				client:   dir.client,
				nodeFlg:  1,
				children: make(map[string]*minioFSNode),
				mode:     0644,
				modTime:  time.Now(),
			}

		}

		if flag&(os.O_SYNC|os.O_APPEND) != 0 {
			// memFile doesn't support these flags yet.
			return nil, os.ErrInvalid
		}
		if flag&os.O_CREATE != 0 {
			if flag&os.O_EXCL != 0 && n != nil {
				return nil, os.ErrExist
			}
			if n == nil {
				n = &minioFSNode{
					mode: perm.Perm(),
				}
				dir.children[frag] = n
			}
		}
		if n == nil {
			return nil, os.ErrNotExist
		}
		if flag&(os.O_WRONLY|os.O_RDWR) != 0 && flag&os.O_TRUNC != 0 {
			n.mu.Lock()
			n.data = nil
			n.mu.Unlock()
		}
	}

	children := make([]os.FileInfo, 0, len(n.children))
	for cName, c := range n.children {
		children = append(children, c.stat(cName))
	}
	return &minioFile{
		n:                n,
		nameSnapshot:     frag,
		childrenSnapshot: children,
	}, nil
}

func (fs *MinioFS) RemoveAll(ctx context.Context, name string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dir, frag, err := fs.find("remove", name)
	if err != nil {
		return err
	}
	if dir == nil {
		// We can't remove the root.
		return os.ErrInvalid
	}
	delete(dir.children, frag)
	return nil
}

func (fs *MinioFS) Rename(ctx context.Context, oldName, newName string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	oldName = slashClean(oldName)
	newName = slashClean(newName)
	if oldName == newName {
		return nil
	}
	if strings.HasPrefix(newName, oldName+"/") {
		// We can't rename oldName to be a sub-directory of itself.
		return os.ErrInvalid
	}

	oDir, oFrag, err := fs.find("rename", oldName)
	if err != nil {
		return err
	}
	if oDir == nil {
		// We can't rename from the root.
		return os.ErrInvalid
	}

	nDir, nFrag, err := fs.find("rename", newName)
	if err != nil {
		return err
	}
	if nDir == nil {
		// We can't rename to the root.
		return os.ErrInvalid
	}

	oNode, ok := oDir.children[oFrag]
	if !ok {
		return os.ErrNotExist
	}
	if oNode.children != nil {
		if nNode, ok := nDir.children[nFrag]; ok {
			if nNode.children == nil {
				return errors.New("webdav: not a directory")
			}
			if len(nNode.children) != 0 {
				return errors.New("webdav: directory not empty")
			}
		}
	}
	delete(oDir.children, oFrag)
	nDir.children[nFrag] = oNode
	return nil
}

func (fs *MinioFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dir, frag, err := fs.find("stat", name)
	if err != nil {
		return nil, err
	}
	if dir == nil {
		// We're stat'ting the root.
		return fs.root.stat("/"), nil
	}
	if n, ok := dir.children[frag]; ok {
		return n.stat(path.Base(name)), nil
	}
	return nil, os.ErrNotExist
}

func slashClean(name string) string {
	if name == "" || name[0] != '/' {
		name = "/" + name
	}
	return path.Clean(name)
}
