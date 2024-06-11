// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	model "deals_chatting_app_backend/internal/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CreateOrUpdatePreferences mocks base method.
func (m *MockUserRepository) CreateOrUpdatePreferences(ctx context.Context, userID uuid.UUID, preferences model.Preferences) (*model.Preferences, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdatePreferences", ctx, userID, preferences)
	ret0, _ := ret[0].(*model.Preferences)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrUpdatePreferences indicates an expected call of CreateOrUpdatePreferences.
func (mr *MockUserRepositoryMockRecorder) CreateOrUpdatePreferences(ctx, userID, preferences interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdatePreferences", reflect.TypeOf((*MockUserRepository)(nil).CreateOrUpdatePreferences), ctx, userID, preferences)
}

// CreateOrUpdateProfile mocks base method.
func (m *MockUserRepository) CreateOrUpdateProfile(ctx context.Context, userID uuid.UUID, profile model.Profile) (*model.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdateProfile", ctx, userID, profile)
	ret0, _ := ret[0].(*model.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrUpdateProfile indicates an expected call of CreateOrUpdateProfile.
func (mr *MockUserRepositoryMockRecorder) CreateOrUpdateProfile(ctx, userID, profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdateProfile", reflect.TypeOf((*MockUserRepository)(nil).CreateOrUpdateProfile), ctx, userID, profile)
}

// FindAll mocks base method.
func (m *MockUserRepository) FindAll(ctx context.Context, userID uuid.UUID) ([]model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, userID)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockUserRepositoryMockRecorder) FindAll(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockUserRepository)(nil).FindAll), ctx, userID)
}

// FindByID mocks base method.
func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockUserRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockUserRepository)(nil).FindByID), ctx, id)
}

// FindByUsername mocks base method.
func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUsername", ctx, username)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUsername indicates an expected call of FindByUsername.
func (mr *MockUserRepositoryMockRecorder) FindByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUsername", reflect.TypeOf((*MockUserRepository)(nil).FindByUsername), ctx, username)
}

// GetProfileByUserID mocks base method.
func (m *MockUserRepository) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileByUserID", ctx, userID)
	ret0, _ := ret[0].(*model.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileByUserID indicates an expected call of GetProfileByUserID.
func (mr *MockUserRepositoryMockRecorder) GetProfileByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileByUserID", reflect.TypeOf((*MockUserRepository)(nil).GetProfileByUserID), ctx, userID)
}

// Save mocks base method.
func (m *MockUserRepository) Save(ctx context.Context, user model.User) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, user)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockUserRepositoryMockRecorder) Save(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserRepository)(nil).Save), ctx, user)
}
