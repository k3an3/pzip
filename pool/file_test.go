package pool_test

import (
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/k3an3/pzip/internal/testutils"
	"github.com/k3an3/pzip/pool"
)

const (
	testdataRoot             = "../testdata/"
	archivePath              = testdataRoot + "archive.zip"
	helloTxtFileFixture      = testdataRoot + "hello.txt"
	helloMarkdownFileFixture = testdataRoot + "hello.md"
	helloDirectoryFixture    = testdataRoot + "hello/"
)

func TestNewFile(t *testing.T) {
	t.Run("with file name relative to archive root when file path is relative", func(t *testing.T) {
		info := testutils.GetFileInfo(t, helloTxtFileFixture)
		file, err := pool.NewFile(helloTxtFileFixture, info, "")
		assert.NoError(t, err)

		assert.Equal(t, "hello.txt", file.Header.Name)
	})

	t.Run("with file name relative to archive root when file path is absolute", func(t *testing.T) {
		absFilePath, err := filepath.Abs(helloTxtFileFixture)
		assert.NoError(t, err)
		info := testutils.GetFileInfo(t, absFilePath)
		file, err := pool.NewFile(absFilePath, info, "")
		assert.NoError(t, err)

		assert.Equal(t, "hello.txt", file.Header.Name)
	})

	t.Run("with file name relative to archive root for directories", func(t *testing.T) {
		filePath := filepath.Join(helloDirectoryFixture, "nested/hello.md")
		info := testutils.GetFileInfo(t, filePath)

		file, err := pool.NewFile(filePath, info, helloDirectoryFixture)
		assert.NoError(t, err)

		assert.Equal(t, "hello/nested/hello.md", file.Header.Name)
	})

	t.Run("resets file as new", func(t *testing.T) {
		filePath := filepath.Join(helloDirectoryFixture, "nested/hello.md")
		info := testutils.GetFileInfo(t, filePath)

		file, err := pool.NewFile(filePath, info, helloDirectoryFixture)
		assert.NoError(t, err)

		newInfo := testutils.GetFileInfo(t, helloTxtFileFixture)
		err = file.Reset(helloTxtFileFixture, newInfo, "")
		assert.NoError(t, err)

		assert.Equal(t, helloTxtFileFixture, file.Path)
		assert.Equal(t, newInfo, file.Info)
		assert.Equal(t, "hello.txt", file.Header.Name)
		assert.Equal(t, 0, file.CompressedData.Len())
		assert.Equal(t, pool.DefaultBufferSize, file.CompressedData.Cap())
	})
}
