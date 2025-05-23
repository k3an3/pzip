package pool_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/k3an3/pzip/pool"
)

func TestFileWorkerPool(t *testing.T) {
	t.Run("can enqueue tasks", func(t *testing.T) {
		fileProcessPool, err := pool.NewFileWorkerPool(func(f *pool.File) error { return nil }, &pool.Config{Concurrency: 1, Capacity: 1})
		assert.NoError(t, err)
		fileProcessPool.Start(context.Background())

		fileProcessPool.Enqueue(&pool.File{})

		assert.Equal(t, 1, fileProcessPool.PendingFiles())
	})

	t.Run("has workers process files to completion", func(t *testing.T) {
		output := bytes.Buffer{}
		executor := func(_ *pool.File) error {
			output.WriteString("hello, world!")
			return nil
		}

		fileProcessPool, err := pool.NewFileWorkerPool(executor, &pool.Config{Concurrency: 1, Capacity: 1})
		assert.NoError(t, err)
		fileProcessPool.Start(context.Background())

		fileProcessPool.Enqueue(&pool.File{})

		err = fileProcessPool.Close()

		assert.NoError(t, err)
		assert.Equal(t, 0, fileProcessPool.PendingFiles())
		assert.Equal(t, "hello, world!", output.String())
	})

	t.Run("returns an error if number of workers is less than one", func(t *testing.T) {
		executor := func(_ *pool.File) error { return nil }

		_, err := pool.NewFileWorkerPool(executor, &pool.Config{Concurrency: 0, Capacity: 1})
		assert.Error(t, err)
	})

	t.Run("can be closed and restarted", func(t *testing.T) {
		output := bytes.Buffer{}
		executor := func(_ *pool.File) error {
			output.WriteString("hello ")
			return nil
		}

		fileProcessPool, err := pool.NewFileWorkerPool(executor, &pool.Config{Concurrency: 1, Capacity: 1})
		assert.NoError(t, err)

		fileProcessPool.Start(context.Background())
		fileProcessPool.Enqueue(&pool.File{})
		err = fileProcessPool.Close()
		assert.NoError(t, err)

		fileProcessPool.Start(context.Background())
		fileProcessPool.Enqueue(&pool.File{})
		err = fileProcessPool.Close()

		assert.NoError(t, err)
		assert.Equal(t, "hello hello ", output.String())
	})

	t.Run("stops workers with first error encountered by a goroutine", func(t *testing.T) {
		executor := func(file *pool.File) error {
			if file.Path == "1" {
				return errors.New("file is corrupt")
			}

			return nil
		}

		fileProcessPool, err := pool.NewFileWorkerPool(executor, &pool.Config{Concurrency: 2, Capacity: 1})
		assert.NoError(t, err)

		fileProcessPool.Start(context.Background())

		fileProcessPool.Enqueue(&pool.File{})
		fileProcessPool.Enqueue(&pool.File{})
		fileProcessPool.Enqueue(&pool.File{Path: "1"})

		err = fileProcessPool.Close()

		assert.Error(t, err)
	})
}
