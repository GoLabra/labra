// Package utils groups helper utilities used throughout the API. This file
// provides an abstraction over the OS file system and resides in src/api/utils.
package utils

import (
	"io"
	"os"
)

// OSFileSystem implements filesystem operations using the OS package.
type OSFileSystem struct{}

// NewOSFileSystem creates an OSFileSystem instance.
func NewOSFileSystem() OSFileSystem {
	return OSFileSystem{}
}

// Open wraps os.Open.
func (OSFileSystem) Open(name string) (*os.File, error) { return os.Open(name) }

// Create wraps os.Create.
func (OSFileSystem) Create(name string) (*os.File, error) { return os.Create(name) }

// Exit terminates the process with the given exit code.
func (OSFileSystem) Exit(code int) { os.Exit(code) }

// Remove wraps os.Remove.
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
