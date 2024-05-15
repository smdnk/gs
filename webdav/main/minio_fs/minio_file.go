package minio_fs

import (
	"encoding/xml"
	"github.com/minio/minio-go/v7"
	"gs/webdav"
	"io"
	"os"
	"time"
)

// minioFile 实现了 webdav.File 接口
type minioFile struct {
	client           *minio.Client
	name             string
	n                *minioFSNode
	nameSnapshot     string
	childrenSnapshot []os.FileInfo
	pos              int
}

func (f *minioFile) DeadProps() (map[xml.Name]webdav.Property, error) { return f.n.DeadProps() }
func (f *minioFile) Patch(patches []webdav.Proppatch) ([]webdav.Propstat, error) {
	return f.n.Patch(patches)
}

func (f *minioFile) Close() error {
	return nil
}

func (f *minioFile) Read(p []byte) (int, error) {
	f.n.mu.Lock()
	defer f.n.mu.Unlock()
	if f.n.mode.IsDir() {
		return 0, os.ErrInvalid
	}
	if f.pos >= len(f.n.data) {
		return 0, io.EOF
	}
	n := copy(p, f.n.data[f.pos:])
	f.pos += n
	return n, nil
}

// Readdir 获取目录的子目录
func (f *minioFile) Readdir(count int) ([]os.FileInfo, error) {
	f.n.mu.Lock()
	defer f.n.mu.Unlock()
	if !f.n.mode.IsDir() {
		return nil, os.ErrInvalid
	}
	old := f.pos
	if old >= len(f.childrenSnapshot) {
		// The os.File Readdir docs say that at the end of a directory,
		// the error is io.EOF if count > 0 and nil if count <= 0.
		if count > 0 {
			return nil, io.EOF
		}
		return nil, nil
	}
	if count > 0 {
		f.pos += count
		if f.pos > len(f.childrenSnapshot) {
			f.pos = len(f.childrenSnapshot)
		}
	} else {
		f.pos = len(f.childrenSnapshot)
		old = 0
	}
	return f.childrenSnapshot[old:f.pos], nil
}

func (f *minioFile) Seek(offset int64, whence int) (int64, error) {
	f.n.mu.Lock()
	defer f.n.mu.Unlock()
	npos := f.pos
	// TODO: How to handle offsets greater than the size of system int?
	switch whence {
	case io.SeekStart:
		npos = int(offset)
	case io.SeekCurrent:
		npos += int(offset)
	case io.SeekEnd:
		npos = len(f.n.data) + int(offset)
	default:
		npos = -1
	}
	if npos < 0 {
		return 0, os.ErrInvalid
	}
	f.pos = npos
	return int64(f.pos), nil
}

func (f *minioFile) Stat() (os.FileInfo, error) {
	return f.n.stat(f.nameSnapshot), nil
}

func (f *minioFile) Write(p []byte) (int, error) {
	lenp := len(p)
	f.n.mu.Lock()
	defer f.n.mu.Unlock()

	if f.n.mode.IsDir() {
		return 0, os.ErrInvalid
	}
	if f.pos < len(f.n.data) {
		n := copy(f.n.data[f.pos:], p)
		f.pos += n
		p = p[n:]
	} else if f.pos > len(f.n.data) {
		// Write permits the creation of holes, if we've seek'ed past the
		// existing end of file.
		if f.pos <= cap(f.n.data) {
			oldLen := len(f.n.data)
			f.n.data = f.n.data[:f.pos]
			hole := f.n.data[oldLen:]
			for i := range hole {
				hole[i] = 0
			}
		} else {
			d := make([]byte, f.pos, f.pos+len(p))
			copy(d, f.n.data)
			f.n.data = d
		}
	}

	if len(p) > 0 {
		// We should only get here if f.pos == len(f.n.data).
		f.n.data = append(f.n.data, p...)
		f.pos = len(f.n.data)
	}
	f.n.modTime = time.Now()
	return lenp, nil
}
