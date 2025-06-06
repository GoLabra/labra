// Code generated by MockGen. DO NOT EDIT.
// Source: ./domain/svc/entity.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/GoLabra/labra/src/api/entgql/entity"
	generator "github.com/GoLabra/labra/src/api/entgql/generator"
	gomock "github.com/golang/mock/gomock"
)

// MockSchemaManager is a mock of SchemaManager interface.
type MockSchemaManager struct {
	ctrl     *gomock.Controller
	recorder *MockSchemaManagerMockRecorder
}

// MockSchemaManagerMockRecorder is the mock recorder for MockSchemaManager.
type MockSchemaManagerMockRecorder struct {
	mock *MockSchemaManager
}

// NewMockSchemaManager creates a new mock instance.
func NewMockSchemaManager(ctrl *gomock.Controller) *MockSchemaManager {
	mock := &MockSchemaManager{ctrl: ctrl}
	mock.recorder = &MockSchemaManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchemaManager) EXPECT() *MockSchemaManagerMockRecorder {
	return m.recorder
}

// BackupSchema mocks base method.
func (m *MockSchemaManager) BackupSchema() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BackupSchema")
	ret0, _ := ret[0].(error)
	return ret0
}

// BackupSchema indicates an expected call of BackupSchema.
func (mr *MockSchemaManagerMockRecorder) BackupSchema() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BackupSchema", reflect.TypeOf((*MockSchemaManager)(nil).BackupSchema))
}

// Generate mocks base method.
func (m *MockSchemaManager) Generate(ctx context.Context, entities []entity.Entity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", ctx, entities)
	ret0, _ := ret[0].(error)
	return ret0
}

// Generate indicates an expected call of Generate.
func (mr *MockSchemaManagerMockRecorder) Generate(ctx, entities interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockSchemaManager)(nil).Generate), ctx, entities)
}

// RemoveEntityFromSchema mocks base method.
func (m *MockSchemaManager) RemoveEntityFromSchema(fileName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveEntityFromSchema", fileName)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveEntityFromSchema indicates an expected call of RemoveEntityFromSchema.
func (mr *MockSchemaManagerMockRecorder) RemoveEntityFromSchema(fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveEntityFromSchema", reflect.TypeOf((*MockSchemaManager)(nil).RemoveEntityFromSchema), fileName)
}

// RevertSchema mocks base method.
func (m *MockSchemaManager) RevertSchema() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevertSchema")
	ret0, _ := ret[0].(error)
	return ret0
}

// RevertSchema indicates an expected call of RevertSchema.
func (mr *MockSchemaManagerMockRecorder) RevertSchema() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevertSchema", reflect.TypeOf((*MockSchemaManager)(nil).RevertSchema))
}

// WriteEntityToSchema mocks base method.
func (m *MockSchemaManager) WriteEntityToSchema(entityTemplateData generator.EntityTemplateData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteEntityToSchema", entityTemplateData)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteEntityToSchema indicates an expected call of WriteEntityToSchema.
func (mr *MockSchemaManagerMockRecorder) WriteEntityToSchema(entityTemplateData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteEntityToSchema", reflect.TypeOf((*MockSchemaManager)(nil).WriteEntityToSchema), entityTemplateData)
}
