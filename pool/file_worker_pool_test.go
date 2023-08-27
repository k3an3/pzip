package pool_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/pzip/internal/testutils"
	"github.com/pzip/pool"
)

const (
	testdataRoot             = "../testdata/"
	archivePath              = testdataRoot + "archive.zip"
	helloTxtFileFixture      = testdataRoot + "hello.txt"
	helloMarkdownFileFixture = testdataRoot + "hello.md"
	helloDirectoryFixture    = testdataRoot + "hello/"
)

func TestFileWorkerPool(t *testing.T) {
	t.Run("can enqueue tasks", func(t *testing.T) {
		fileProcessPool, err := pool.NewFileWorkerPool(1, func(f pool.File) {})
		assert.NoError(t, err)

		info := testutils.GetFileInfo(t, helloTxtFileFixture)
		fileProcessPool.Enqueue(pool.File{Path: helloTxtFileFixture, Info: info})

		assert.Equal(t, 1, fileProcessPool.PendingFiles())
	})

	t.Run("has workers process files to completion", func(t *testing.T) {
		output := bytes.Buffer{}
		executor := func(_ pool.File) {
			time.Sleep(5 * time.Millisecond)
			output.WriteString("hello, world!")
		}

		fileProcessPool, err := pool.NewFileWorkerPool(1, executor)
		assert.NoError(t, err)
		fileProcessPool.Start()

		info := testutils.GetFileInfo(t, helloTxtFileFixture)
		fileProcessPool.Enqueue(pool.File{Path: helloTxtFileFixture, Info: info})

		fileProcessPool.Close()

		assert.Equal(t, 0, fileProcessPool.PendingFiles())
		assert.Equal(t, "hello, world!", output.String())
	})

	t.Run("returns an error if number of workers is less than one", func(t *testing.T) {
		executor := func(_ pool.File) {
		}
		_, err := pool.NewFileWorkerPool(0, executor)
		assert.Error(t, err)
	})

	t.Run("can be closed and restarted", func(t *testing.T) {
		output := bytes.Buffer{}
		executor := func(_ pool.File) {
			output.WriteString("hello ")
		}

		fileProcessPool, err := pool.NewFileWorkerPool(1, executor)
		assert.NoError(t, err)
		fileProcessPool.Start()

		info := testutils.GetFileInfo(t, helloTxtFileFixture)
		fileProcessPool.Enqueue(pool.File{Path: helloTxtFileFixture, Info: info})

		fileProcessPool.Close()

		fileProcessPool.Start()
		info = testutils.GetFileInfo(t, helloTxtFileFixture)
		fileProcessPool.Enqueue(pool.File{Path: helloTxtFileFixture, Info: info})

		fileProcessPool.Close()

		assert.Equal(t, "hello hello ", output.String())
	})
}
