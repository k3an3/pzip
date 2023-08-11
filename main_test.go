package main

import (
	"archive/zip"
	"fmt"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
)

const (
	archivePath = "testdata/archive.zip"
	srcRoot     = "testdata/"
)

func TestArchive(t *testing.T) {
	t.Run("archives a single non-empty file with a name", func(t *testing.T) {
		file := openTestFile(t, srcRoot+"hello.txt")
		defer file.Close()

		archive, cleanup := createTempArchive(t, archivePath)
		defer cleanup()

		archiver := NewArchiver(archive)
		archiver.Archive(file)

		archiveReader := getArchiveReader(t, archive.Name())
		defer archiveReader.Close()

		assert.Equal(t, 1, len(archiveReader.File))
		assertArchiveContainsFile(t, archiveReader.File, "hello.txt")

		info, err := file.Stat()
		assert.NoError(t, err)

		got := archiveReader.File[0].UncompressedSize64
		want := uint64(info.Size())

		assert.Equal(t, want, got, "expected file %s to have raw size %d but got %d", file.Name(), want, got)
	})

	t.Run("archives two files sequentially", func(t *testing.T) {
		file1 := openTestFile(t, srcRoot+"hello.txt")
		defer file1.Close()
		file2 := openTestFile(t, srcRoot+"/hello.md")
		defer file2.Close()

		archive, cleanup := createTempArchive(t, archivePath)
		defer cleanup()

		archiver := NewArchiver(archive)
		archiver.Archive(file1, file2)

		archiveReader := getArchiveReader(t, archive.Name())
		defer archiveReader.Close()

		assert.Equal(t, 2, len(archiveReader.File))
	})
}

func openTestFile(t testing.TB, name string) *os.File {
	t.Helper()

	file, err := os.Open(name)
	assert.NoError(t, err, fmt.Sprintf("could not open %s: %v", name, err))

	return file
}

func createTempArchive(t testing.TB, name string) (*os.File, func()) {
	t.Helper()

	archive, err := os.Create(name)
	assert.NoError(t, err, fmt.Sprintf("could not create archive %s: %v", name, err))

	cleanup := func() {
		archive.Close()
		os.RemoveAll(archive.Name())
	}

	return archive, cleanup
}

func getArchiveReader(t testing.TB, name string) *zip.ReadCloser {
	t.Helper()

	reader, err := zip.OpenReader(name)
	assert.NoError(t, err)

	return reader
}

func assertArchiveContainsFile(t testing.TB, files []*zip.File, name string) {
	t.Helper()

	for _, f := range files {
		if f.Name == name {
			return
		}
	}

	t.Errorf("expected file %s to be in archive but wasn't", name)
}
