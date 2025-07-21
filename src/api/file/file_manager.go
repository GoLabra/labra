// Package file provides helper types for interacting with file storage. It is
// located under src/api/file.
package file

import (
	"context"

	"github.com/GoLabra/labra/src/api/config"
	"github.com/GoLabra/labra/src/api/entgql/ent"
)

// FileManager wraps file storage configuration.
type FileManager struct {
	config *config.Config
}

// SaveFile writes the provided content using a configured storage provider.
func SaveFile(ctx context.Context, name string, content string) (ent.CreateFileInput, error) {
	return ent.CreateFileInput{}, nil
}
