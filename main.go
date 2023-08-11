package main

import (
	"archive/zip"
	"io"
	"os"
)

type Archiver struct {
	Dest *os.File
	w    *zip.Writer
}

func NewArchiver(archive *os.File) *Archiver {
	return &Archiver{Dest: archive, w: zip.NewWriter(archive)}
}

func (a *Archiver) Archive(files ...*os.File) error {
	for _, file := range files {
		a.writeFile(file)
	}

	a.w.Close()
	return nil
}

func (a *Archiver) writeFile(file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	writer, err := a.w.Create(info.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}

func main() {
}
