package coresize

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"os"
	"path"

	_ "image/jpeg"
	_ "image/png"
)

const filechunk = 8192 // 8k

type ImageFile struct {
	Path string
	Hash string
}

func (i ImageFile) Name() string {
	return path.Base(i.Path)
}

func (i ImageFile) NameWithHash() string {
	return fmt.Sprintf("%s-%s", i.Hash, i.Name())
}

func (i *ImageFile) ComputeHash() error {
	file, err := os.Open(i.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf))
	}

	i.Hash = fmt.Sprintf("%x", hash.Sum(nil))[:8]
	return nil
}