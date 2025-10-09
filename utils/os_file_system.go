package utils

import (
	"io"
	"os"
)

type OSFileSystem struct {
}

func NewOSFileSystem() OSFileSystem {
	return OSFileSystem{}
}

func (OSFileSystem) Open(name string) (*os.File, error) { return os.Open(name) }

func (OSFileSystem) Create(name string) (*os.File, error) { return os.Create(name) }

func (OSFileSystem) Exit(code int) { os.Exit(code) }

func (OSFileSystem) Remove(name string) error { return os.Remove(name) }

func (fs OSFileSystem) CopyDirectory(from, to string) error {
	err := os.RemoveAll(to)
	if err != nil {
		return err
	}

	err = os.MkdirAll(to, os.ModePerm)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(from)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fromPath := from + string(os.PathSeparator) + entry.Name()
		toPath := to + string(os.PathSeparator) + entry.Name()

		if entry.IsDir() {
			err = fs.CopyDirectory(fromPath, toPath)
			if err != nil {
				return err
			}
			continue
		}

		toF, err := os.Create(toPath)
		if err != nil {
			return err
		}
		defer toF.Close()

		fromF, err := os.Open(fromPath)
		if err != nil {
			return err
		}
		defer fromF.Close()

		_, err = io.Copy(toF, fromF)
		if err != nil {
			return err
		}
	}

	return nil
}
