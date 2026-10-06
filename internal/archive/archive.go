package archive

import (
	"archive/tar"
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
)

type Archive struct {
	tar  *tar.Writer
	zstd *zstd.Encoder
}

func NewArchive(w io.Writer) (*Archive, error) {
	zstdEncoder, err := zstd.NewWriter(w)
	if err != nil {
		return nil, err
	}

	tarWriter := tar.NewWriter(zstdEncoder)

	return &Archive{
		tar:  tarWriter,
		zstd: zstdEncoder,
	}, nil
}

func (a *Archive) Close() error {
	if err := a.tar.Close(); err != nil {
		return err
	}
	if err := a.zstd.Close(); err != nil {
		return err
	}
	return nil
}

func (a *Archive) AddFile(name string, file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name: name,
		Mode: 0600,
		Size: info.Size(),
	}

	if err := a.tar.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(a.tar, file); err != nil {
		return err
	}

	return nil
}
