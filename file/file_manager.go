package file

import (
	"context"

	"github.com/GoLabra/labra/config"
	"github.com/GoLabra/labra/entgql/ent"
)

type FileManager struct {
	config *config.Config
}

func SaveFile(ctx context.Context, name string, content string) (ent.CreateFileInput, error) {
	return ent.CreateFileInput{}, nil
}
