package svc

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoLabra/labra/src/api/config"
	"github.com/GoLabra/labra/src/api/entgql/domain/repo"
	"github.com/GoLabra/labra/src/api/entgql/ent"
	"github.com/samborkent/uuidv7"
)

type FileManager interface {
	Save(ctx context.Context, fileId string) (*ent.File, error)
	Get(ctx context.Context, fileId string) (*ent.File, error)
}

type File struct {
	fileManager FileManager
	repository  *repo.Repository
}

func NewFile(r *repo.Repository) *File {
	return &File{repository: r}
}

func (s *File) Get(ctx context.Context, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) ([]*ent.File, error) {
	config, ok := ctx.Value("config").(*config.Config)

	if !ok {
		return nil, fmt.Errorf("[File.Get] config not set in context")
	}

	files, err := s.repository.File.Get(ctx, where, orderBy, skip, first, last)
	if err != nil {
		return files, err
	}

	for idx, file := range files {
		fileLocation := filepath.Join(config.FileStoragePath, file.StorageFileName)
		fileContent, err := os.ReadFile(fileLocation)
		if err != nil {
			return files, err
		}

		files[idx].Content = base64.StdEncoding.EncodeToString(fileContent)
	}

	return files, err
}

func (s *File) Connection(ctx context.Context, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) (*ent.FileConnection, error) {
	return s.repository.File.Connection(ctx, where, orderBy, skip, first, last)
}

func (s *File) GetOne(ctx context.Context, where ent.FileWhereUniqueInput) (*ent.File, error) {
	config, ok := ctx.Value("config").(*config.Config)

	if !ok {
		return nil, fmt.Errorf("[File.Get] config not set in context")
	}

	file, err := s.repository.File.GetOne(ctx, where)
	if err != nil {
		return file, err
	}

	fileLocation := filepath.Join(config.FileStoragePath, file.StorageFileName)
	fileContent, err := os.ReadFile(fileLocation)
	if err != nil {
		return file, err
	}

	file.Content = base64.StdEncoding.EncodeToString(fileContent)

	return file, err
}

func (s *File) GetTx(ctx context.Context, tx *ent.Tx, where *ent.FileWhereInput, orderBy *ent.FileOrder, skip *int, first *int, last *int) ([]*ent.File, error) {
	return s.repository.File.GetTx(ctx, tx, where, orderBy, skip, first, last)
}

func (s *File) GetOneTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput) (*ent.File, error) {
	return s.repository.File.GetOneTx(ctx, tx, where)
}
func (s *File) CreateTx(ctx context.Context, tx *ent.Tx, data ent.CreateFileInput) (*ent.File, error) {
	config, ok := ctx.Value("config").(*config.Config)

	if !ok {
		return nil, fmt.Errorf("[File.CreateTx] config not set in context")
	}

	fileExtension := filepath.Ext(data.Name)

	if data.Caption == nil {
		caption := strings.TrimSuffix(data.Name, fileExtension)
		data.Caption = &caption
	}

	data.StorageFileName = uuidv7.New().String() + "." + fileExtension

	decoded, err := base64.StdEncoding.DecodeString(data.Content)
	if err != nil {
		return nil, fmt.Errorf("[File.CreateTx] failed to decode base64 string: %w", err)
	}

	data.Size = int64(len(decoded))

	outputPath := filepath.Join(config.FileStoragePath, data.StorageFileName)

	err = os.WriteFile(outputPath, decoded, 0644)
	if err != nil {
		return nil, fmt.Errorf("[File.CreateTx] failed to write file: %w", err)
	}

	createdInput, err := s.repository.File.CreateTx(ctx, tx, data)
	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (s *File) Create(ctx context.Context, data ent.CreateFileInput) (*ent.File, error) {
	config, ok := ctx.Value("config").(*config.Config)

	if !ok {
		return nil, fmt.Errorf("[File.CreateTx] config not set in context")
	}

	fileExtension := filepath.Ext(data.Name)

	if data.Caption == nil {
		caption := strings.TrimSuffix(data.Name, fileExtension)
		data.Caption = &caption
	}

	data.StorageFileName = uuidv7.New().String() + "." + fileExtension

	decoded, err := base64.StdEncoding.DecodeString(data.Content)
	if err != nil {
		return nil, fmt.Errorf("[File.CreateTx] failed to decode base64 string: %w", err)
	}

	data.Size = int64(len(decoded))

	outputPath := filepath.Join(config.FileStoragePath, data.StorageFileName)

	err = os.WriteFile(outputPath, decoded, 0644)
	if err != nil {
		return nil, fmt.Errorf("[File.CreateTx] failed to write file: %w", err)
	}

	createdInput, err := s.repository.File.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return createdInput, err
}

func (s *File) CreateMany(ctx context.Context, data []ent.CreateFileInput) ([]*ent.File, error) {
	return s.repository.File.CreateMany(ctx, data)
}

func (s *File) CreateManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateFileInput) ([]*ent.File, error) {
	return s.repository.File.CreateManyTx(ctx, tx, data)
}

func (s *File) Update(ctx context.Context, where ent.FileWhereUniqueInput, data ent.UpdateFileInput) (*ent.File, error) {
	return s.repository.File.Update(ctx, where, data)
}

func (s *File) UpdateTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput, data ent.UpdateFileInput) (*ent.File, error) {
	return s.repository.File.UpdateTx(ctx, tx, where, data)
}

func (s *File) UpdateMany(ctx context.Context, where ent.FileWhereInput, data ent.UpdateFileInput) (int, error) {
	return s.repository.File.UpdateMany(ctx, where, data)
}

func (s *File) UpdateManyTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereInput, data ent.UpdateFileInput) (int, error) {
	return s.repository.File.UpdateManyTx(ctx, tx, where, data)
}

func (s *File) Upsert(ctx context.Context, data ent.CreateFileInput) (*ent.File, error) {
	return s.repository.File.Upsert(ctx, data)
}

func (s *File) UpsertTx(ctx context.Context, tx *ent.Tx, data ent.CreateFileInput) (*ent.File, error) {
	return s.repository.File.UpsertTx(ctx, tx, data)
}

func (s *File) UpsertMany(ctx context.Context, data []ent.CreateFileInput) (int, error) {
	return s.repository.File.UpsertMany(ctx, data)
}

func (s *File) UpsertManyTx(ctx context.Context, tx *ent.Tx, data []ent.CreateFileInput) (int, error) {
	return s.repository.File.UpsertManyTx(ctx, tx, data)
}

func (s *File) Delete(ctx context.Context, where ent.FileWhereUniqueInput) (*ent.File, error) {
	return s.repository.File.Delete(ctx, where)
}

func (s *File) DeleteTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereUniqueInput) (*ent.File, error) {
	return s.repository.File.DeleteTx(ctx, tx, where)
}

func (s *File) DeleteMany(ctx context.Context, where ent.FileWhereInput) (int, error) {
	return s.repository.File.DeleteMany(ctx, where)
}

func (s *File) DeleteManyTx(ctx context.Context, tx *ent.Tx, where ent.FileWhereInput) (int, error) {
	return s.repository.File.DeleteManyTx(ctx, tx, where)
}
